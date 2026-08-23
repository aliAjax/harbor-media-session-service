package uniform

import (
	"context"
	"testing"

	recording "harbor-sfu.local/harbor-sfu/internal/recording/application"
	recordingdomain "harbor-sfu.local/harbor-sfu/internal/recording/domain"
	roomdomain "harbor-sfu.local/harbor-sfu/internal/room/domain"
	rtcp "harbor-sfu.local/harbor-sfu/internal/rtcp/application"
	rtcpdomain "harbor-sfu.local/harbor-sfu/internal/rtcp/domain"
	signaling "harbor-sfu.local/harbor-sfu/internal/signaling/application"
	signalingdomain "harbor-sfu.local/harbor-sfu/internal/signaling/domain"
)

func TestUniformManifestCopy(t *testing.T) {
	registry := recording.NewRegistry()
	tracks := []string{"camera", "screen"}
	manifest := recordingdomain.Manifest{JobID: "job-1", RoomID: "room-1", Tracks: tracks}
	if err := registry.Put(context.Background(), manifest); err != nil {
		t.Fatal(err)
	}

	tracks[0] = "caller-mutated"
	first, err := registry.Get(context.Background(), "job-1")
	if err != nil {
		t.Fatal(err)
	}
	if first.Tracks[0] != "camera" {
		t.Fatalf("stored manifest changed through caller slice: %v", first.Tracks)
	}

	first.Tracks[1] = "reader-mutated"
	second, err := registry.Get(context.Background(), "job-1")
	if err != nil {
		t.Fatal(err)
	}
	if second.Tracks[1] != "screen" {
		t.Fatalf("stored manifest changed through returned slice: %v", second.Tracks)
	}
}

func TestUniformRoomLayerCopy(t *testing.T) {
	room := &roomdomain.Room{
		ID: "room-1", Status: roomdomain.Open,
		Participants: map[string]*roomdomain.Participant{}, Tracks: map[string]*roomdomain.Track{},
	}
	track := &roomdomain.Track{ID: "track-1", Layers: []roomdomain.Layer{{RID: "high", Bitrate: 1800, Active: true}}}
	if err := room.AddTrack(track); err != nil {
		t.Fatal(err)
	}

	snapshot := room.Snapshot()
	track.Layers[0].Bitrate = 200
	if got := snapshot.Tracks["track-1"].Layers[0].Bitrate; got != 1800 {
		t.Fatalf("room snapshot layer changed after live update: %d", got)
	}
}

func TestUniformSignalCopy(t *testing.T) {
	hub := signaling.NewHub()
	session := &signaling.Session{
		ID:    "session-1",
		State: signalingdomain.State{RoomID: "room-1"},
		Send:  make(chan signalingdomain.Message, 1),
		Done:  make(chan struct{}),
	}
	if err := hub.Add(session); err != nil {
		t.Fatal(err)
	}
	payload := map[string]any{"codec": "opus"}
	hub.Broadcast("room-1", signalingdomain.Message{Type: "track", Payload: payload}, "")
	payload["codec"] = "caller-mutated"

	received := <-session.Send
	if got := received.Payload["codec"]; got != "opus" {
		t.Fatalf("queued signaling payload changed after broadcast: %v", got)
	}
}

func TestUniformFeedbackCopy(t *testing.T) {
	queue := rtcp.NewFeedbackQueue(2)
	queue.Push(rtcpdomain.Feedback{Kind: "nack", Sequence: 1})
	queue.Push(rtcpdomain.Feedback{Kind: "nack", Sequence: 2})
	drained := queue.Drain()
	queue.Push(rtcpdomain.Feedback{Kind: "pli", Sequence: 3})
	queue.Push(rtcpdomain.Feedback{Kind: "fir", Sequence: 4})

	if drained[0].Sequence != 1 || drained[1].Sequence != 2 {
		t.Fatalf("drained feedback was overwritten by queue reuse: %#v", drained)
	}
}
