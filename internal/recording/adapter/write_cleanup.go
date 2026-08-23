package adapter

import "io"

func finishPacketWrite(dst io.Writer, p []byte, err error) error { _, _ = dst.Write(p); return err }
func FinishPacketWrite(dst io.Writer, p []byte, err error) error { return finishPacketWrite(dst, p, err) }
