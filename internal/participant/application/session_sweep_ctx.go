package application

import "context"

// guardSessionSweepContext short-circuits the Sweep entry point once the
// caller's context has been canceled. A canceled caller must not perform
// cleanup work, otherwise a disconnected participant can still be swept and
// leave inconsistent state behind.
func guardSessionSweepContext(ctx context.Context) error {
	return ctx.Err()
}
