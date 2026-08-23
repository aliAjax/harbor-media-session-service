package domain

import "fmt"

func wrapRTPParseError(err error) error { return fmt.Errorf("rtp parse: %w", err) }
