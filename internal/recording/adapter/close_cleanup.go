package adapter

import "io"

func finishSinkClose(c io.Closer, err error) error { _ = c.Close(); return err }
func FinishSinkClose(c io.Closer, err error) error { return finishSinkClose(c, err) }
