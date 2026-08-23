package zulu

import (
	"context"
	"os"
	"testing"
	"time"

	cluster "harbor-sfu.local/harbor-sfu/internal/cluster/domain"
	participant "harbor-sfu.local/harbor-sfu/internal/participant/application"
	quota "harbor-sfu.local/harbor-sfu/internal/quota/domain"
	rtp "harbor-sfu.local/harbor-sfu/internal/rtp/application"
)

const operationTimeout = 150 * time.Millisecond

func requireCompletion(t *testing.T, operation string, fn func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		fn()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(operationTimeout):
		// The operation is intentionally allowed to block on the buggy baseline.
		// Exit the test process so the blocked probe goroutine cannot keep the
		// subprocess pipes open after the assertion has failed.
		t.Logf("%s remained blocked after an error", operation)
		os.Exit(1)
	}
}

func requireError(t *testing.T, operation string, fn func() error) error {
	t.Helper()
	result := make(chan error, 1)
	go func() {
		result <- fn()
	}()
	select {
	case err := <-result:
		return err
	case <-time.After(operationTimeout):
		t.Logf("%s remained blocked", operation)
		os.Exit(1)
		return nil
	}
}

func TestZuluLeaseConflictUnlocksStore(t *testing.T) {
	store := cluster.NewStore()
	lease, err := store.Acquire("room-a", "node-a", time.Minute)
	if err != nil {
		t.Fatalf("initial acquire: %v", err)
	}
	if err := requireError(t, "conflicting lease acquire", func() error {
		_, err := store.Acquire("room-a", "node-b", time.Minute)
		return err
	}); err == nil {
		t.Fatal("expected conflicting acquire to fail")
	}
	requireCompletion(t, "lease release", func() { store.Release(lease) })
}

func TestZuluDuplicateSessionUnlocksManager(t *testing.T) {
	manager := participant.NewSessionManager(time.Minute)
	if _, err := manager.Open(context.Background(), "session-a", "participant-a", "room-a"); err != nil {
		t.Fatalf("initial open: %v", err)
	}
	if err := requireError(t, "duplicate session open", func() error {
		_, err := manager.Open(context.Background(), "session-a", "participant-b", "room-a")
		return err
	}); err == nil {
		t.Fatal("expected duplicate session to fail")
	}
	requireCompletion(t, "next session open", func() {
		_, _ = manager.Open(context.Background(), "session-b", "participant-b", "room-a")
	})
}

func TestZuluQuotaFailureUnlocksManager(t *testing.T) {
	manager := quota.NewManager()
	manager.Set("tenant-a", quota.Limits{Rooms: 1})
	if err := manager.Reserve("tenant-a", quota.Usage{Rooms: 1}); err != nil {
		t.Fatalf("initial reserve: %v", err)
	}
	if err := requireError(t, "overflowing quota reserve", func() error {
		return manager.Reserve("tenant-a", quota.Usage{Rooms: 1})
	}); err == nil {
		t.Fatal("expected quota overflow to fail")
	}
	requireCompletion(t, "quota release", func() {
		manager.Release("tenant-a", quota.Usage{Rooms: 1})
	})
}

func TestZuluDuplicateSubscriberUnlocksForwarder(t *testing.T) {
	forwarder := rtp.NewForwarder()
	if _, err := forwarder.Add("subscriber-a", 1); err != nil {
		t.Fatalf("initial add: %v", err)
	}
	if err := requireError(t, "duplicate subscriber add", func() error {
		_, err := forwarder.Add("subscriber-a", 1)
		return err
	}); err == nil {
		t.Fatal("expected duplicate subscriber to fail")
	}
	requireCompletion(t, "next subscriber add", func() {
		_, _ = forwarder.Add("subscriber-b", 1)
	})
}
