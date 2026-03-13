package windows

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/CedricThomas/console/internal/boundary/out/command"
)

// executor implements the CommandExecutor interface for Windows systems
type executor struct{}

// New creates a new Windows command executor
func New() command.CommandExecutor {
	return &executor{}
}

// Shutdown powers off the Windows system using shutdown /p command
func (e *executor) Shutdown(ctx context.Context) error {
	_, err := e.Execute(ctx, "shutdown", "/p")
	if err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	return nil
}
