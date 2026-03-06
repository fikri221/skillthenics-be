package worker

import "context"

// Worker defines a common interface for background services.
type Worker interface {
	// Start begins the worker's execution. It should be non-blocking.
	Start(ctx context.Context)
}
