package domain_test

import (
	"context"
	"errors"
	clusterapp "harbor-sfu.local/harbor-sfu/internal/cluster/application"
	clusterdomain "harbor-sfu.local/harbor-sfu/internal/cluster/domain"
	"testing"
	"time"
)

func canceledPapaContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}
func TestPapaCanceledClaim(t *testing.T) {
	s := clusterdomain.NewStore()
	c := clusterapp.NewCoordinator(s, "node-a", time.Minute)
	if _, err := c.Claim(canceledPapaContext(), "room-a"); !errors.Is(err, context.Canceled) {
		t.Fatalf("claim cancellation lost: %v", err)
	}
	if _, err := s.Acquire("room-a", "node-b", time.Minute); err != nil {
		t.Fatalf("canceled claim changed lease state: %v", err)
	}
}
func TestPapaCanceledRenew(t *testing.T) {
	s := clusterdomain.NewStore()
	l, err := s.Acquire("room-b", "node-a", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	c := clusterapp.NewCoordinator(s, "node-a", time.Hour)
	if _, err = c.Renew(canceledPapaContext(), l); !errors.Is(err, context.Canceled) {
		t.Fatalf("renew cancellation lost: %v", err)
	}
}
func TestPapaExpiredLeaseRejected(t *testing.T) {
	s := clusterdomain.NewStore()
	l, err := s.Acquire("room-c", "node-a", -time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.Renew(l, time.Minute); !errors.Is(err, clusterdomain.ErrLeaseLost) {
		t.Fatalf("expired lease renewed: %v", err)
	}
}
func TestPapaForgedOwnerCannotRelease(t *testing.T) {
	s := clusterdomain.NewStore()
	l, err := s.Acquire("room-d", "node-a", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	l.Owner = "node-b"
	s.Release(l)
	if _, err = s.Acquire("room-d", "node-c", time.Minute); !errors.Is(err, clusterdomain.ErrLeaseLost) {
		t.Fatalf("forged owner released lease: %v", err)
	}
}
