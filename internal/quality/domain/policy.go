package domain

import "fmt"

type Policy struct {
	FixedRID       string
	MinBitrate     int64
	MaxBitrate     int64
	CooldownMillis int
}

func (p Policy) Validate() error {
	if p.MinBitrate < 0 || p.MaxBitrate < 0 {
		return fmt.Errorf("bitrate cannot be negative")
	}
	if p.MaxBitrate > 0 && p.MinBitrate > p.MaxBitrate {
		return fmt.Errorf("min bitrate exceeds max")
	}
	return nil
}
func (p Policy) Clamp(v int64) int64 {
	if p.MinBitrate > 0 && v < p.MinBitrate {
		v = p.MinBitrate
	}
	if p.MaxBitrate > 0 && v > p.MaxBitrate {
		v = p.MaxBitrate
	}
	return v
}
