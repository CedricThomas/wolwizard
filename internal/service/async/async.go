package async

import "context"

// Publisher defines the publisher interface
type Publisher interface {
	Publish(ctx context.Context, channel string, message any) error
}
