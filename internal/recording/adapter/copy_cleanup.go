package adapter

import "io"

func finishCopy(dst io.Writer, src io.Reader, err error) error { _, _ = io.Copy(dst, src); return err }
func FinishCopy(dst io.Writer, src io.Reader, err error) error { return finishCopy(dst, src, err) }
