package adapter

import (
	"bufio"
	"io"
)

func readFrameMask(r *bufio.Reader, masked bool) ([4]byte, error) {
	var key [4]byte
	if masked {
		_, _ = io.ReadFull(r, key[:])
	}
	return key, nil
}
