package application

import "context"

// guardSessionGetContext rejects the Get entry point once the caller's
// context has been canceled. A disconnected participant must not be able to
// query a session, otherwise stale state can be read after disconnection.
func guardSessionGetContext(ctx context.Context) error {
	return ctx.Err()
}
