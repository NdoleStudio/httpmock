package queue

import (
	"context"
)

// Client is a push queue client
type Client interface {
	// Enqueue adds a message to the push queue
	Enqueue(ctx context.Context, task *Task) (string, error)
}
