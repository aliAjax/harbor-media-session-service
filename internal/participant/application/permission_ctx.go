package application

import "context"

// guardPermissionContext rejects the permission-check entry points once the
// caller's context has been canceled. A disconnected participant must not be
// allowed to proceed through a permission decision, otherwise authorization
// keeps running after the connection is already gone.
func guardPermissionContext(ctx context.Context) error {
	return ctx.Err()
}
