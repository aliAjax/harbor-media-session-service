package cinder

import (
	"errors"
	quality "harbor-sfu.local/harbor-sfu/internal/quality/application"
	qualitydomain "harbor-sfu.local/harbor-sfu/internal/quality/domain"
	recording "harbor-sfu.local/harbor-sfu/internal/recording/domain"
	rtcp "harbor-sfu.local/harbor-sfu/internal/rtcp/adapter"
	rtp "harbor-sfu.local/harbor-sfu/internal/rtp/domain"
	"testing"
)

func runWithoutPanic(t *testing.T, name string, fn func()) {
	t.Helper()
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("%s panicked: %v", name, recovered)
		}
	}()
	fn()
}

func TestCinderEmptyManifestReturnsError(t *testing.T) {
	runWithoutPanic(t, "empty manifest validation", func() {
		manifest := recording.Manifest{JobID: "job-cinder", RoomID: "room-cinder"}
		if err := manifest.Validate(); err == nil {
			t.Fatal("empty manifest unexpectedly passed validation")
		}
	})
}

func TestCinderNoLayerUsesFallback(t *testing.T) {
	runWithoutPanic(t, "empty quality decision", func() {
		decision := quality.Decide(qualitydomain.Policy{}, nil, 0, 64000)
		if decision.RID != "low" || decision.Bitrate != 64000 {
			t.Fatalf("empty quality decision = %#v, want low fallback", decision)
		}
	})
}

func TestCinderShortRTPReturnsError(t *testing.T) {
	runWithoutPanic(t, "short RTP parsing", func() {
		if _, err := rtp.Parse([]byte{0x80}); !errors.Is(err, rtp.ErrShort) {
			t.Fatalf("short RTP error = %v, want ErrShort", err)
		}
	})
}

func TestCinderShortRTCPReturnsError(t *testing.T) {
	runWithoutPanic(t, "short RTCP parsing", func() {
		if _, err := rtcp.ParseReport([]byte{0x80, 200, 0, 4}); err == nil {
			t.Fatal("short RTCP report unexpectedly parsed")
		}
	})
}
