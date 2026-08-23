package application

import (
	"context"
	"harbor-sfu.local/harbor-sfu/internal/rtp/domain"
	"net"
)

func validateRunInputs(_ context.Context, _ *net.UDPConn, _ func(domain.Packet)) error { return nil }
