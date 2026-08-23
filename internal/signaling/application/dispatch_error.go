package application

import (
	"errors"
	"fmt"
)

func wrapDispatchError(err error) error { return fmt.Errorf("dispatch: %w", err) }

var ErrStaleMessage = errors.New("stale signaling sequence")
