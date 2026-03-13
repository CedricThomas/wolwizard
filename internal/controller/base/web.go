package base

import (
	"context"
	"fmt"

	"github.com/CedricThomas/console/internal/controller"
	"github.com/CedricThomas/console/internal/domain"
	asyncapi "github.com/CedricThomas/console/internal/input/async/api"
	"github.com/CedricThomas/console/internal/service/async"
	"github.com/CedricThomas/console/internal/service/keystore"
)

type web struct {
	publisher async.Publisher
	keystore  keystore.Keystore
}

func NewWebController(publisher async.Publisher, keystore keystore.Keystore) controller.Web {
	return &web{
		publisher: publisher,
		keystore:  keystore,
	}
}

func (w web) SendAsyncBootCommand(ctx context.Context, osName domain.OSName) error {
	// Publish boot command to Redis pubsub channel
	bootCmd := asyncapi.BootCommand{OSName: osName}
	if err := w.publisher.Publish(ctx, asyncapi.BootChannel, bootCmd); err != nil {
		return fmt.Errorf("publish boot command: %w", err)
	}

	return nil
}

func (w web) SendAsyncShutdownCommand(ctx context.Context) error {
	// Publish shutdown command to Redis pubsub channel
	shutdownCmd := asyncapi.ShutdownCommand{}
	if err := w.publisher.Publish(ctx, asyncapi.ShutdownChannel, shutdownCmd); err != nil {
		return fmt.Errorf("publish shutdown command: %w", err)
	}

	return nil
}

func (w web) ProcessMetrics(_ context.Context, _ domain.Metrics) error {
	// TODO implement
	return nil
}
