package domain

type LayerInput struct {
	RID     string
	Bitrate int64
	Active  bool
}
type LayerDecision struct {
	RID     string
	Bitrate int64
	Reason  string
	Changed bool
}
type Metrics struct {
	Packets   uint64
	Bytes     uint64
	Lost      uint64
	Jitter    uint32
	RTTMillis int64
}
