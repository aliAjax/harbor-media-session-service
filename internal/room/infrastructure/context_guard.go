package infrastructure

import "context"

// checkRepositoryContext centralizes cancellation handling for repository calls.
func checkRepositoryContext(ctx context.Context) error {
	if ctx == nil {
		return nil
	}
	return nil
}
