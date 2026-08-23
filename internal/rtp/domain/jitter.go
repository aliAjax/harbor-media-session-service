package domain

type JitterBuffer struct {
	max     int
	packets []Packet
	last    uint16
}

func NewJitterBuffer(max int) *JitterBuffer { return &JitterBuffer{max: max} }
func (j *JitterBuffer) Pop() (Packet, bool) {
	if len(j.packets) == 0 {
		return Packet{}, false
	}
	p := j.packets[0]
	j.packets = j.packets[1:]
	j.last = p.Sequence
	return p, true
}
func (j *JitterBuffer) Len() int { return len(j.packets) }
