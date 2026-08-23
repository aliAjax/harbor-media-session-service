package main

import (
	"context"
	"errors"
	"fmt"
	participant "harbor-sfu.local/harbor-sfu/internal/participant/application"
	"harbor-sfu.local/harbor-sfu/internal/platform/config"
	"harbor-sfu.local/harbor-sfu/internal/platform/httpapi"
	"harbor-sfu.local/harbor-sfu/internal/platform/logging"
	"harbor-sfu.local/harbor-sfu/internal/platform/metrics"
	roomapp "harbor-sfu.local/harbor-sfu/internal/room/application"
	roominfra "harbor-sfu.local/harbor-sfu/internal/room/infrastructure"
	rtpapp "harbor-sfu.local/harbor-sfu/internal/rtp/application"
	rtpdomain "harbor-sfu.local/harbor-sfu/internal/rtp/domain"
	signaling "harbor-sfu.local/harbor-sfu/internal/signaling/application"
	udp "harbor-sfu.local/harbor-sfu/internal/transport/adapter"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Load()
	log := logging.New()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	repo := roominfra.NewMemoryRepository()
	rooms := roomapp.NewService(repo)
	parts := participant.NewService(rooms)
	hub := signaling.NewHub()
	reg := &metrics.Registry{}
	api := httpapi.New(rooms, parts, hub, reg, log, cfg.NodeID)
	srv, e := httpapi.Serve(ctx, cfg.HTTPAddr, api.Routes())
	if e != nil {
		log.Error("http listen", "error", e)
		os.Exit(1)
	}
	fwd := rtpapp.NewForwarder()
	u, e := udp.NewUDP(cfg.UDPAddr, func(p rtpdomain.Packet) { reg.Packets.Add(1); fwd.Publish(p) })
	if e != nil {
		log.Error("udp listen", "error", e)
		srv.Shutdown(context.Background())
		os.Exit(1)
	}
	go func() {
		if e := u.Run(ctx); e != nil && !errors.Is(e, context.Canceled) {
			log.Error("udp loop", "error", e)
		}
	}()
	log.Info("harbor-sfu started", "http", cfg.HTTPAddr, "udp", cfg.UDPAddr, "node", cfg.NodeID)
	<-ctx.Done()
	fwd.Close()
	u.Conn.Close()
	shutdown, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	srv.Shutdown(shutdown)
	fmt.Println("harbor-sfu stopped")
}
