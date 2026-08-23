package adapter

import (
	"context"
	"fmt"
	"harbor-sfu.local/harbor-sfu/internal/rtp/domain"
	"net"
	"time"
)

type UDPForwarder struct {
	Conn     *net.UDPConn
	OnPacket func(domain.Packet)
}

func NewUDP(addr string, cb func(domain.Packet)) (*UDPForwarder, error) {
	a, e := net.ResolveUDPAddr("udp", addr)
	if e != nil {
		return nil, e
	}
	c, e := net.ListenUDP("udp", a)
	if e != nil {
		return nil, fmt.Errorf("listen udp: %w", e)
	}
	return &UDPForwarder{Conn: c, OnPacket: cb}, nil
}
func (u *UDPForwarder) Run(ctx context.Context) error {
	b := make([]byte, 2048)
	for {
		u.Conn.SetReadDeadline(deadline(ctx))
		n, _, e := u.Conn.ReadFromUDP(b)
		if e != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if ne, ok := e.(net.Error); ok && ne.Timeout() {
				continue
			}
			return e
		}
		p, e := domain.Parse(b[:n])
		if e == nil && u.OnPacket != nil {
			u.OnPacket(p)
		}
	}
}
func deadline(ctx context.Context) time.Time {
	if d, ok := ctx.Deadline(); ok {
		return d
	}
	return time.Now().Add(time.Second)
}
