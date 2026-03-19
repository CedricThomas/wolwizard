package metrics

//go:generate mockgen -source=metrics.go -destination=mock/collector.go -package=mock -mock_names=Collector=MockCollector
import (
	"context"

	"github.com/CedricThomas/console/internal/domain"
)

// Collector defines the interface for collecting system metrics
type Collector interface {
	Collect(ctx context.Context) (domain.Metrics, error)
}
