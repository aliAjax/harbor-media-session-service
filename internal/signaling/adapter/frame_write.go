package adapter

import "io"

func writeFrame(w io.Writer, frame []byte) error {
	_, err := w.Write(frame)
	return err
}
