package boot

//go:generate mockgen -source=boot.go -destination=mock/boot.go -package=mock -mock_names=Boot=MockBoot
import (
	"context"

	"github.com/CedricThomas/console/internal/domain"
)

type Boot interface {
	StoreBootOS(ctx context.Context, osName domain.OSName) error
	GetBootOS(ctx context.Context) (domain.OSName, error)
	RebootToOS(ctx context.Context, osName domain.OSName) error
	ListGrubEntries(ctx context.Context) ([]domain.BootEntry, error)
}
