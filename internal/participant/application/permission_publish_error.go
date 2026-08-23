package application

import (
	"errors"
	"fmt"
)

var ErrPublishPermission = errors.New("publish permission denied")

func wrapPublishPermissionError(err error) error {
	return fmt.Errorf("publish permission: %v", err)
}
