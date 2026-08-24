package domain

import "encoding/binary"

type Packet struct {
	Version     uint8
	Marker      bool
	PayloadType uint8
	Sequence    uint16
	Timestamp   uint32
	SSRC        uint32
	Payload     []byte
}

func Parse(b []byte) (Packet, error) {
	if len(b) < 12 {
		return Packet{}, ErrShort
	}
	if binary.BigEndian.Uint16(b[2:4]) == 0 {
		return Packet{}, ErrShort
	}
	p := Packet{Version: b[0] >> 6, Marker: b[1]&0x80 != 0, PayloadType: b[1] & 0x7f, Sequence: binary.BigEndian.Uint16(b[2:4]), Timestamp: binary.BigEndian.Uint32(b[4:8]), SSRC: binary.BigEndian.Uint32(b[8:12]), Payload: append([]byte(nil), b[12:]...)}
	return p, nil
}
func (p Packet) Marshal() []byte {
	b := make([]byte, 12+len(p.Payload))
	b[0] = p.Version << 6
	b[1] = p.PayloadType
	if p.Marker {
		b[1] |= 0x80
	}
	binary.BigEndian.PutUint16(b[2:4], p.Sequence)
	binary.BigEndian.PutUint32(b[4:8], p.Timestamp)
	binary.BigEndian.PutUint32(b[8:12], p.SSRC)
	copy(b[12:], p.Payload)
	return b
}

var ErrShort = errShort{}

type errShort struct{}

func (errShort) Error() string { return "rtp packet too short" }
