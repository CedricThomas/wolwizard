package metrics

import (
	"context"

	"github.com/CedricThomas/console/internal/domain"
)

type Metrics interface {
	ProcessMetrics(ctx context.Context, metrics domain.Metrics) error
	GetLastMetrics(ctx context.Context) (*domain.Metrics, error)
}
