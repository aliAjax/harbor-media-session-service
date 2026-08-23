package httpapi

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"harbor-sfu.local/harbor-sfu/internal/participant/application"
	"harbor-sfu.local/harbor-sfu/internal/platform/metrics"
	roomapp "harbor-sfu.local/harbor-sfu/internal/room/application"
	"harbor-sfu.local/harbor-sfu/internal/room/domain"
	"harbor-sfu.local/harbor-sfu/internal/signaling/adapter"
	signaling "harbor-sfu.local/harbor-sfu/internal/signaling/application"
	signdomain "harbor-sfu.local/harbor-sfu/internal/signaling/domain"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

type Server struct {
	Rooms        *roomapp.Service
	Participants *application.Service
	Hub          *signaling.Hub
	Metrics      *metrics.Registry
	Logger       *slog.Logger
	NodeID       string
	closed       atomic.Bool
}

func New(r *roomapp.Service, p *application.Service, h *signaling.Hub, m *metrics.Registry, l *slog.Logger, node string) *Server {
	return &Server{Rooms: r, Participants: p, Hub: h, Metrics: m, Logger: l, NodeID: node}
}
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.health)
	mux.HandleFunc("/readyz", s.ready)
	mux.HandleFunc("/metrics", s.Metrics.Handler)
	mux.HandleFunc("/v1/rooms", s.rooms)
	mux.HandleFunc("/v1/rooms/", s.roomAction)
	return s.middleware(mux)
}
func (s *Server) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Metrics.Requests.Add(1)
		w.Header().Set("X-Request-ID", fmt.Sprintf("req-%d", time.Now().UnixNano()))
		defer func() {
			if v := recover(); v != nil {
				http.Error(w, "internal error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	write(w, 200, map[string]any{"status": "ok", "node": s.NodeID})
}
func (s *Server) ready(w http.ResponseWriter, _ *http.Request) {
	if s.closed.Load() {
		write(w, 503, map[string]string{"status": "draining"})
		return
	}
	write(w, 200, map[string]string{"status": "ready"})
}
func (s *Server) rooms(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tenant := r.URL.Query().Get("tenant")
		v, e := s.Rooms.List(r.Context(), tenant)
		if e != nil {
			errorJSON(w, e)
			return
		}
		write(w, 200, v)
	case http.MethodPost:
		var in struct {
			ID       string `json:"id"`
			TenantID string `json:"tenant_id"`
			Name     string `json:"name"`
		}
		if json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&in) != nil || in.ID == "" {
			errorJSON(w, errors.New("id is required"))
			return
		}
		v, e := s.Rooms.Create(r.Context(), in.TenantID, in.ID, in.Name, s.NodeID)
		if e != nil {
			errorJSON(w, e)
			return
		}
		write(w, 201, v)
	default:
		http.Error(w, "method not allowed", 405)
	}
}
func (s *Server) roomAction(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		http.NotFound(w, r)
		return
	}
	id := parts[2]
	if len(parts) == 3 {
		v, e := s.Rooms.Get(r.Context(), id)
		if e != nil {
			errorJSON(w, e)
			return
		}
		write(w, 200, v)
		return
	}
	action := parts[3]
	switch action {
	case "join":
		var in struct{ ID, Identity, Role string }
		json.NewDecoder(r.Body).Decode(&in)
		if in.ID == "" {
			in.ID = "participant-" + fmt.Sprint(time.Now().UnixNano())
		}
		p, e := s.Participants.Join(r.Context(), id, in.ID, in.Identity, in.Role)
		if e != nil {
			errorJSON(w, e)
			return
		}
		write(w, 201, p)
	case "publish":
		var in struct {
			ID, Kind, Codec, PublisherID string
			SSRC                         uint32
			Layers                       []domain.Layer
		}
		json.NewDecoder(r.Body).Decode(&in)
		e := s.Participants.Publish(r.Context(), id, &domain.Track{ID: in.ID, Kind: in.Kind, Codec: in.Codec, PublisherID: in.PublisherID, SSRC: in.SSRC, Layers: in.Layers, CreatedAt: time.Now()})
		if e != nil {
			errorJSON(w, e)
			return
		}
		write(w, 201, map[string]string{"status": "published"})
	case "close", "drain":
		st := domain.Closed
		if action == "drain" {
			st = domain.Draining
		}
		if e := s.Rooms.Transition(r.Context(), id, st); e != nil {
			errorJSON(w, e)
			return
		}
		write(w, 200, map[string]string{"status": string(st)})
	case "signal":
		s.websocket(w, r, id)
	default:
		http.NotFound(w, r)
	}
}
func (s *Server) websocket(w http.ResponseWriter, r *http.Request, roomID string) {
	h, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "websocket unsupported", 500)
		return
	}
	c, br, e := h.Hijack()
	if e != nil {
		return
	}
	if e = adapter.Upgrade(c, br.Reader, c); e != nil {
		c.Close()
		return
	}
	sid := "ws-" + fmt.Sprint(time.Now().UnixNano())
	sess := &signaling.Session{ID: sid, State: signdomain.State{Connected: true, RoomID: roomID}, Send: make(chan signdomain.Message, 32), Done: make(chan struct{})}
	s.Hub.Add(sess)
	defer s.Hub.Remove(sid)
	go func() {
		for m := range sess.Send {
			adapter.WriteMessage(c, m)
		}
	}()
	for {
		m, e := adapter.ReadMessage(br.Reader)
		if e != nil {
			break
		}
		s.Metrics.Messages.Add(1)
		if !sess.State.Accept(m) {
			continue
		}
		if m.Type == "leave" {
			break
		}
		m.RoomID = roomID
		m.ParticipantID = sess.State.ParticipantID
		s.Hub.Broadcast(roomID, m, sid)
		adapter.WriteMessage(c, signdomain.Message{Type: "ack", ID: m.ID, Seq: m.Seq, RoomID: roomID})
	}
	c.Close()
}
func write(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
func errorJSON(w http.ResponseWriter, e error) {
	status := 500
	if errors.Is(e, domain.ErrNotFound) {
		status = 404
	}
	if errors.Is(e, domain.ErrConflict) {
		status = 409
	}
	write(w, status, map[string]string{"error": e.Error()})
}
func Serve(ctx context.Context, addr string, h http.Handler) (*http.Server, error) {
	srv := &http.Server{Addr: addr, Handler: h, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second}
	ln, e := net.Listen("tcp", addr)
	if e != nil {
		return nil, e
	}
	go func() { <-ctx.Done(); srv.Shutdown(context.Background()); ln.Close() }()
	go srv.Serve(ln)
	return srv, nil
}
func _unused(*bufio.Reader) {}
