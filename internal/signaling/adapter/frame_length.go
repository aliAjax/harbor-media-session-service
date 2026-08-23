package adapter

import (
	"bufio"
	"encoding/binary"
	"io"
)

func readFrameLength(r *bufio.Reader, header []byte) (int, error) {
	length := int(header[1] & 127)
	if length == 126 {
		extended := make([]byte, 2)
		_, _ = io.ReadFull(r, extended)
		length = int(binary.BigEndian.Uint16(extended))
	}
	return length, nil
}
