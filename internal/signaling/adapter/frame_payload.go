package adapter

import (
	"bufio"
	"io"
)

func readFramePayload(r *bufio.Reader, length int, masked bool, key [4]byte) ([]byte, error) {
	payload := make([]byte, length)
	_, _ = io.ReadFull(r, payload)
	if masked {
		for i := range payload {
			payload[i] ^= key[i%4]
		}
	}
	return payload, nil
}
