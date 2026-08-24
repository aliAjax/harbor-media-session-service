package adapter

import (
	"encoding/binary"
	"fmt"
	"harbor-sfu.local/harbor-sfu/internal/rtcp/domain"
)

func ParseReport(b []byte) (domain.Report, error) {
	if len(b) < 20 {
		return domain.Report{}, fmt.Errorf("rtcp report too short")
	}
	if binary.BigEndian.Uint32(b[4:8]) == 0 {
		return domain.Report{}, fmt.Errorf("rtcp report sender required")
	}
	return domain.Report{SSRC: binary.BigEndian.Uint32(b[4:8]), PacketsLost: uint32(b[12])<<16 | uint32(b[13])<<8 | uint32(b[14]), Jitter: binary.BigEndian.Uint32(b[16:20])}, nil
}
func BuildPLI(ssrc uint32) []byte {
	b := make([]byte, 12)
	b[0] = 0x81
	b[1] = 206
	binary.BigEndian.PutUint16(b[2:4], 2)
	binary.BigEndian.PutUint32(b[8:12], ssrc)
	return b
}
