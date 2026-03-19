package async

//go:generate mockgen -source=async.go -destination=mock/publisher.go -package=mock -mock_names=Publisher=MockPublisher
import "context"

// Publisher defines the publisher interface
type Publisher interface {
	Publish(ctx context.Context, channel string, message any) error
}
