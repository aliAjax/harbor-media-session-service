package application

import "context"

// guardSessionOpenContext rejects the Open entry point once the caller's
// context has been canceled. A disconnected participant must not be able to
// open a new session, otherwise the session lingers as dirty state.
func guardSessionOpenContext(ctx context.Context) error {
	return ctx.Err()
}
