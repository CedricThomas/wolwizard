package async

// Publisher defines the Redis publisher interface
type Publisher interface {
	Publish(channel string, message []byte) error
	Close() error
}
