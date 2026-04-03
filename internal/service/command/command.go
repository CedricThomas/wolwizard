package command

//go:generate mockgen -source=command.go -destination=mock/command.go -package=mock -mock_names=CommandExecutor=MockCommandExecutor
import (
	"context"

	"github.com/CedricThomas/console/internal/domain"
)

// CommandExecutor defines the interface for executing system commands
type CommandExecutor interface {
	Shutdown(ctx context.Context) error
	SetGrubReboot(ctx context.Context, entryName string) error
	Reboot(ctx context.Context) error
	ListGrubEntries(ctx context.Context) ([]domain.BootEntry, error)
}

type ErrUnsupportedOS struct{}

func (e *ErrUnsupportedOS) Error() string {
	return "operation not supported on this OS"
}
