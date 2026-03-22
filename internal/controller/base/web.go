package base

import (
	"context"
	"fmt"
	"time"

	"github.com/CedricThomas/console/internal/config"
	"github.com/CedricThomas/console/internal/controller"
	"github.com/CedricThomas/console/internal/domain"
	asyncapi "github.com/CedricThomas/console/internal/input/async/api"
	"github.com/CedricThomas/console/internal/service/async"
	"github.com/CedricThomas/console/internal/service/keystore"
	"github.com/CedricThomas/console/internal/service/token"
	"github.com/CedricThomas/console/internal/service/websocket"
	metricsusecase "github.com/CedricThomas/console/internal/usecase/metrics"
	metricsusecasebase "github.com/CedricThomas/console/internal/usecase/metrics/base"
)

type web struct {
	auth
	publisher      async.Publisher
	keystore       keystore.Keystore
	metricsUsecase metricsusecase.Metrics
}

func NewWebController(publisher async.Publisher, keystore keystore.Keystore, tokenSrv token.Service, cfg *config.WebConfig, wsManager websocket.Manager) controller.Web {
	authCtrl := newAuthController(keystore, tokenSrv)
	metricsUsecase := metricsusecasebase.New(keystore, time.Duration(cfg.LastMetricsKeyTTLSeconds)*time.Second, wsManager)
	return &web{
		auth:           authCtrl,
		publisher:      publisher,
		keystore:       keystore,
		metricsUsecase: metricsUsecase,
	}
}

// Web controller methods
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

func (w web) ProcessMetrics(ctx context.Context, metrics domain.Metrics) error {
	if err := w.metricsUsecase.ProcessMetrics(ctx, metrics); err != nil {
		return fmt.Errorf("process metrics: %w", err)
	}
	return nil
}
