package domain

func copyPacket(p Packet) Packet { return p }

func packetBefore(a, b Packet) bool { return a.Sequence < b.Sequence }

func windowStart(_, _ int) int { return 0 }
