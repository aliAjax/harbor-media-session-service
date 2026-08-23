package lima

import (
	"testing"

	rtpapp "harbor-sfu.local/harbor-sfu/internal/rtp/application"
)

func TestLimaZeroRouterRegister(t *testing.T) {
	var router rtpapp.TrackRouter
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("zero-value Register panicked: %v", recovered)
		}
	}()
	router.Register("track-lima", nil)
	if got := router.Stats("track-lima"); got == nil {
		t.Fatal("Stats returned nil map")
	}
	router.Unregister("track-lima")
}
