package adapter

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"harbor-sfu.local/harbor-sfu/internal/signaling/domain"
	"io"
	"net"
	"strings"
)

type Conn struct {
	net.Conn
	Reader *bufio.Reader
}

func Upgrade(w io.Writer, r *bufio.Reader, c net.Conn) error {
	line, e := r.ReadString('\n')
	if e != nil {
		return e
	}
	_ = line
	key := ""
	for {
		line, e = r.ReadString('\n')
		if e != nil {
			return e
		}
		if line == "\r\n" {
			break
		}
		if len(line) > 5 && line[:5] == "Sec-W" {
			parts := strings.SplitN(strings.TrimSpace(line), ":", 2)
			if len(parts) == 2 {
				key = strings.TrimSpace(parts[1])
			}
		}
	}
	sum := sha1.Sum([]byte(key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	fmt.Fprintf(w, "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: %s\r\n\r\n", base64.StdEncoding.EncodeToString(sum[:]))
	return nil
}
func ReadMessage(r *bufio.Reader) (domain.Message, error) {
	h := make([]byte, 2)
	if _, e := io.ReadFull(r, h); e != nil {
		return domain.Message{}, e
	}
	n, err := readFrameLength(r, h)
	if err != nil {
		return domain.Message{}, err
	}
	if n > 1<<20 {
		return domain.Message{}, fmt.Errorf("message too large")
	}
	mask := h[1]&128 != 0
	key, err := readFrameMask(r, mask)
	if err != nil {
		return domain.Message{}, err
	}
	b, err := readFramePayload(r, n, mask, key)
	if err != nil {
		return domain.Message{}, err
	}
	var m domain.Message
	if e := json.Unmarshal(b, &m); e != nil {
		return m, e
	}
	return m, nil
}
func WriteMessage(w io.Writer, m domain.Message) error {
	b, e := json.Marshal(m)
	if e != nil {
		return e
	}
	if len(b) > 125 {
		return fmt.Errorf("message too large")
	}
	return writeFrame(w, append([]byte{0x81, byte(len(b))}, b...))
}
