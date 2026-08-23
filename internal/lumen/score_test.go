package lumen_test

import (
	"context"
	"errors"
	"testing"

	app "harbor-sfu.local/harbor-sfu/internal/room/application"
	domain "harbor-sfu.local/harbor-sfu/internal/room/domain"
	infra "harbor-sfu.local/harbor-sfu/internal/room/infrastructure"
)

func roomForTest(id string) *domain.Room {
	return &domain.Room{ID: id, TenantID: "tenant", Status: domain.Open,
		Participants: map[string]*domain.Participant{}, Tracks: map[string]*domain.Track{}}
}

func TestLumenRepositoryCreateCanceled(t *testing.T) {
	repo := infra.NewMemoryRepository()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := repo.Create(ctx, roomForTest("cancel-create")); !errors.Is(err, context.Canceled) {
		t.Fatalf("Create error = %v, want context.Canceled", err)
	}
	if _, err := repo.Get(context.Background(), "cancel-create"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("canceled Create changed state: Get error = %v", err)
	}
	if err := repo.Create(context.Background(), roomForTest("live-create")); err != nil {
		t.Fatalf("live Create error = %v", err)
	}
}

func TestLumenRepositoryGetCanceled(t *testing.T) {
	repo := infra.NewMemoryRepository()
	if err := repo.Create(context.Background(), roomForTest("stored")); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if room, err := repo.Get(ctx, "stored"); !errors.Is(err, context.Canceled) || room != nil {
		t.Fatalf("Get canceled = (%v, %v), want (nil, context.Canceled)", room, err)
	}
	if room, err := repo.Get(context.Background(), "stored"); err != nil || room == nil {
		t.Fatalf("live Get = (%v, %v), want stored room", room, err)
	}
}

type recordingRepo struct {
	room       *domain.Room
	createCall int
	saveCall   int
}

func (r *recordingRepo) Create(context.Context, *domain.Room) error { r.createCall++; return nil }
func (r *recordingRepo) Get(context.Context, string) (*domain.Room, error) {
	if r.room == nil {
		return nil, domain.ErrNotFound
	}
	return r.room, nil
}
func (r *recordingRepo) List(context.Context, string) ([]domain.Room, error) { return nil, nil }
func (r *recordingRepo) Save(context.Context, *domain.Room) error            { r.saveCall++; return nil }

func TestLumenServiceCreateCanceled(t *testing.T) {
	repo := &recordingRepo{}
	svc := app.NewService(repo)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if room, err := svc.Create(ctx, "tenant", "service-create", "room", "node"); !errors.Is(err, context.Canceled) || room != nil {
		t.Fatalf("Create canceled = (%v, %v), want (nil, context.Canceled)", room, err)
	}
	if repo.createCall != 0 {
		t.Fatalf("canceled Create called repository %d times", repo.createCall)
	}
}

func TestLumenServiceTransitionCanceled(t *testing.T) {
	repo := &recordingRepo{room: roomForTest("service-transition")}
	svc := app.NewService(repo)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := svc.Transition(ctx, "service-transition", domain.Closed); !errors.Is(err, context.Canceled) {
		t.Fatalf("Transition canceled error = %v, want context.Canceled", err)
	}
	if repo.saveCall != 0 || repo.room.Status != domain.Open {
		t.Fatalf("canceled Transition changed state: save=%d status=%s", repo.saveCall, repo.room.Status)
	}
}
