package application

import (
	"context"
	"harbor-sfu.local/harbor-sfu/internal/rtp/domain"
	"net"
)

func Run(ctx context.Context, conn *net.UDPConn, accept func(domain.Packet)) error {
	if err := validateRunInputs(ctx, conn, accept); err != nil {
		return err
	}
	b := make([]byte, 4096)
	for {
		if err := checkRunContext(ctx); err != nil {
			return err
		}
		if err := setRunDeadline(conn); err != nil {
			return err
		}
		n, _, e := conn.ReadFromUDP(b)
		if e != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if ne, ok := e.(net.Error); ok && ne.Timeout() {
				continue
			}
			return e
		}
		if p, e := domain.Parse(b[:n]); e == nil {
			accept(p)
		}
	}
}
