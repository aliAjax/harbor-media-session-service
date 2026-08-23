package domain

type Report struct {
	SSRC        uint32
	PacketsLost uint32
	Jitter      uint32
	RTTMillis   int64
	Bandwidth   int64
}
type Feedback struct {
	Kind     string
	SSRC     uint32
	Sequence uint16
}
