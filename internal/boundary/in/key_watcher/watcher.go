package key_watcher

// Watcher defines a keystore keys watcher interface
type Watcher interface {
	Watch(key string, handler redis.KeyExpirationHandler) error
	Unwatch(key string) error
	Start()
	Stop()
}
