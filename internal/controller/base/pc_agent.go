package base

import (
	"context"
	"fmt"

	"github.com/CedricThomas/console/internal/controller"
	"github.com/CedricThomas/console/internal/domain"
	"github.com/CedricThomas/console/internal/input/async/api"
	"github.com/CedricThomas/console/internal/input/async/presenters"
	"github.com/CedricThomas/console/internal/service/async"
	"github.com/CedricThomas/console/internal/service/command"
)

type pcAgent struct {
	executor  command.CommandExecutor
	publisher async.Publisher
}

func NewPCAgentController(executor command.CommandExecutor) controller.PCAgent {
	return &pcAgent{
		executor: executor,
	}
}

func (pa *pcAgent) ShutdownCurrentHost(ctx context.Context) error {
	err := pa.executor.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("shutdown current host: %w", err)
	}
	return nil
}

// SendAsyncMetrics sends the metrics asynchronously.
func (pa *pcAgent) SendAsyncMetrics(ctx context.Context, metrics domain.Metrics) error {
	err := pa.publisher.Publish(ctx, api.MetricsChannel, presenters.DomainToMetricsCommand(metrics))
	if err != nil {
		return fmt.Errorf("send async metrics: %w", err)
	}
	return nil
}
