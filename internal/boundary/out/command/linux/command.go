package linux

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/CedricThomas/console/internal/boundary/out/command"
)

// executor implements the CommandExecutor interface for Linux systems
type executor struct{}

// New creates a new Linux command executor
func New() command.CommandExecutor {
	return &executor{}
}

// Shutdown powers off the Linux system using poweroff command
func (e *executor) Shutdown(ctx context.Context) error {
	_, err := e.Execute(ctx, "poweroff")
	if err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	return nil
}
