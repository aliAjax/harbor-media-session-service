package tango

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	r "harbor-sfu.local/harbor-sfu/internal/rtp/domain"
	a "harbor-sfu.local/harbor-sfu/internal/signaling/adapter"
	s "harbor-sfu.local/harbor-sfu/internal/signaling/application"
	sd "harbor-sfu.local/harbor-sfu/internal/signaling/domain"
	tr "harbor-sfu.local/harbor-sfu/internal/transport/application"
	"io"
	"net"
	"testing"
)

func TestTangoRTPWrappedCause(t *testing.T) {
	_, e := r.Parse([]byte{1})
	if !errors.Is(e, r.ErrShort) {
		t.Fatal(e)
	}
}
func TestTangoFrameWrappedCause(t *testing.T) {
	_, e := a.ReadMessage(bufio.NewReader(bytes.NewReader([]byte{0x81, 5, '{', '"'})))
	if !errors.Is(e, io.ErrUnexpectedEOF) {
		t.Fatal(e)
	}
}
func TestTangoDispatchSentinel(t *testing.T) {
	h := s.NewHub()
	x := &s.Session{State: sd.State{LastSeq: 2}, Send: make(chan sd.Message, 1), Done: make(chan struct{})}
	e := s.NewDispatcher(h, 0).Dispatch(context.Background(), x, sd.Message{Seq: 1})
	if !errors.Is(e, s.ErrStaleMessage) {
		t.Fatal(e)
	}
}
func TestTangoUDPConcreteCause(t *testing.T) {
	c, _ := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1")})
	c.Close()
	e := tr.Run(context.Background(), c, func(r.Packet) {})
	var o *net.OpError
	if !errors.As(e, &o) {
		t.Fatal(e)
	}
}
