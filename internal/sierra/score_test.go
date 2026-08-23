package sierra

import (
	d "harbor-sfu.local/harbor-sfu/internal/room/domain"
	"testing"
)

func TestSierraActiveTransitions(t *testing.T) {
	if !d.AllowedTransition(d.Open, d.Draining) || !d.AllowedTransition(d.Open, d.Closed) {
		t.Fatal()
	}
}
func TestSierraClosedIsTerminal(t *testing.T) {
	if d.AllowedTransition(d.Closed, d.Open) || d.AllowedTransition(d.Closed, d.Draining) {
		t.Fatal()
	}
}
func TestSierraDrainingTransitions(t *testing.T) {
	if !d.AllowedTransition(d.Draining, d.Closed) || d.AllowedTransition(d.Draining, d.Open) {
		t.Fatal()
	}
}
func TestSierraReopenRejected(t *testing.T) {
	if d.AllowedTransition(d.Closed, d.Open) {
		t.Fatal()
	}
}
