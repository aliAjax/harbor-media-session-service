package domain

func SequenceDistance(a, b uint16) int32 { return int32(int16(a - b)) }
func IsNewer(a, b uint16) bool           { return SequenceDistance(a, b) > 0 }
func LostBetween(last, current uint16) uint16 {
	d := SequenceDistance(current, last)
	if d <= 1 {
		return 0
	}
	return uint16(d - 1)
}
func TimestampDistance(a, b uint32) int64 { return int64(int32(a - b)) }
