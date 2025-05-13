package tokenbucket

import "context"

// Bucket tries to acquire a token if available, or blocks.
type Bucket interface {
	Take(ctx context.Context) error
	TryTake(ctx context.Context) error
}
