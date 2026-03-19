package command

//go:generate mockgen -source=command.go -destination=mock/command.go -package=mock -mock_names=CommandExecutor=MockCommandExecutor,PlatformExecutor=MockPlatformExecutor
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
