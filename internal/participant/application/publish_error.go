package application

import "fmt"

func wrapTrackPublishError(err error) error {
	return fmt.Errorf("publish room lookup: %v", err)
}
