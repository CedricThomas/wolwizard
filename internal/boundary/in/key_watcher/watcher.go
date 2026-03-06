package key_watcher

import "context"

// Watcher defines a keystore keys watcher interface
type Watcher interface {
	Watch(key string, handler KeyWatcher) error
	Unwatch(key string) error
	Start()
	Stop()
}

type KeyWatcher func(ctx context.Context, value string) error
