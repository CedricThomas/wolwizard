package metrics

//go:generate mockgen -source=metrics.go -destination=mock/metrics.go -package=mock -mock_names=Metrics=MockMetrics
import (
	"context"

	"github.com/CedricThomas/console/internal/domain"
)

type Metrics interface {
	ProcessMetrics(ctx context.Context, metrics domain.Metrics) error
	GetLastMetrics(ctx context.Context) (*domain.Metrics, error)
}
