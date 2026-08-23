package november

import (
	"context"
	"errors"
	app "harbor-sfu.local/harbor-sfu/internal/participant/application"
	"testing"
	"time"
)

func c() context.Context { x, c := context.WithCancel(context.Background()); c(); return x }
func TestNovemberCanceledPermission(t *testing.T) {
	if !errors.Is(app.NewPermissionService().CheckPublish(c(), "publisher"), context.Canceled) {
		t.Fatal("permission")
	}
}
func TestNovemberCanceledSessionGet(t *testing.T) {
	m := app.NewSessionManager(time.Hour)
	m.Open(context.Background(), "s", "p", "r")
	if _, e := m.Get(c(), "s"); !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
}
func TestNovemberCanceledSessionOpen(t *testing.T) {
	m := app.NewSessionManager(time.Hour)
	if _, e := m.Open(c(), "s", "p", "r"); !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
}
func TestNovemberCanceledSessionSweep(t *testing.T) {
	m := app.NewSessionManager(-time.Second)
	m.Open(context.Background(), "s", "p", "r")
	if n := m.Sweep(c()); n != 0 {
		t.Fatal(n)
	}
}
