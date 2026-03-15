package base

import (
	"context"
	"fmt"
	"log"

	"github.com/CedricThomas/console/internal/controller"
	"github.com/CedricThomas/console/internal/input/async/api"
	"github.com/CedricThomas/console/internal/input/async/presenters"
	"github.com/CedricThomas/console/internal/service/async"
	"github.com/CedricThomas/console/internal/service/command"
	"github.com/CedricThomas/console/internal/service/metrics"
)

type pcAgent struct {
	executor  command.CommandExecutor
	publisher async.Publisher
	collector metrics.Collector
	auth      controller.Auth
}

func NewPCAgentController(executor command.CommandExecutor, collector metrics.Collector, publisher async.Publisher, authCtrl controller.Auth) controller.PCAgent {
	return &pcAgent{
		executor:  executor,
		collector: collector,
		publisher: publisher,
		auth:      authCtrl,
	}
}

func (pa *pcAgent) CreateAccount(ctx context.Context, username, password string) error {
	return pa.auth.CreateAccount(ctx, username, password)
}

func (pa *pcAgent) ShutdownCurrentHost(ctx context.Context) error {
	err := pa.executor.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("shutdown current host: %w", err)
	}
	return nil
}

// SendAsyncMetrics sends the metrics asynchronously.
func (pa *pcAgent) SendCurrentHostAsyncMetrics(ctx context.Context) error {
	log.Println("Sending current host async metrics")
	hostMetrics, err := pa.collector.Collect(ctx)
	if err != nil {
		return fmt.Errorf("collect metrics: %w", err)
	}
	err = pa.publisher.Publish(ctx, api.MetricsChannel, presenters.DomainToMetricsCommand(hostMetrics))
	if err != nil {
		return fmt.Errorf("send async metrics: %w", err)
	}
	return nil
}
