package windows

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/CedricThomas/console/internal/domain"
	"github.com/CedricThomas/console/internal/service/command"
)

// executor implements the CommandExecutor interface for Windows systems
type executor struct{}

// New creates a new Windows command executor
func New() command.CommandExecutor {
	return &executor{}
}

// Shutdown powers off the Windows system using shutdown /p command
func (e *executor) Shutdown(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "shutdown", "/p")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	return nil
}

// SetGrubReboot does not support GRUB on Windows
func (e *executor) SetGrubReboot(ctx context.Context, entryName string) error {
	return &command.ErrUnsupportedOS{}
}

// Reboot does not use GRUB on Windows
func (e *executor) Reboot(ctx context.Context) error {
	return &command.ErrUnsupportedOS{}
}

// ListGrubEntries does not support GRUB on Windows
func (e *executor) ListGrubEntries(ctx context.Context) ([]domain.BootEntry, error) {
	return nil, &command.ErrUnsupportedOS{}
}
