package async

// Publisher defines the publisher interface
type Publisher interface {
	Publish(channel string, message []byte) error
	Close() error
}
