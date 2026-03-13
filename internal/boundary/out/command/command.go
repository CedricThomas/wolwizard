package command

import (
	"context"
)

// CommandExecutor defines the interface for executing system commands
type CommandExecutor interface {
	// Shutdown powers off the system
	Shutdown(ctx context.Context) error
}

// PlatformExecutor is a concrete executor for a specific platform
type PlatformExecutor interface {
	CommandExecutor
	GetPlatformName() string
}
