package linux

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/CedricThomas/console/internal/domain"
	"github.com/CedricThomas/console/internal/service/command"
)

// executor implements the CommandExecutor interface for Linux systems
type executor struct{}

// New creates a new Linux command executor
func New() command.CommandExecutor {
	return &executor{}
}

// Shutdown powers off the Linux system using poweroff command
func (e *executor) Shutdown(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "poweroff")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	return nil
}

// SetGrubReboot sets the next boot to the specified GRUB entry
func (e *executor) SetGrubReboot(ctx context.Context, entryName string) error {
	// Use single quotes to handle special characters in entry name
	cmd := exec.CommandContext(ctx, "grub-reboot", "'"+entryName+"'")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("set grub reboot: %w", err)
	}
	return nil
}

// Reboot immediately reboots the system
func (e *executor) Reboot(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "reboot")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("reboot: %w", err)
	}
	return nil
}

// ListGrubEntries parses /boot/grub/grub.cfg and returns all menu entries
func (e *executor) ListGrubEntries(ctx context.Context) ([]domain.BootEntry, error) {
	file, err := os.Open("/boot/grub/grub.cfg")
	if err != nil {
		return nil, fmt.Errorf("open grub.cfg: %w", err)
	}
	defer file.Close()

	var entries []domain.BootEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Match lines starting with "menuentry '
		if strings.HasPrefix(line, "menuentry '") {
			// Extract the entry name between the single quotes
			parts := strings.Split(line, "'")
			if len(parts) > 1 {
				entryName := parts[1]
				if entryName != "" {
					entries = append(entries, domain.BootEntry{Name: entryName})
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read grub.cfg: %w", err)
	}

	return entries, nil
}
