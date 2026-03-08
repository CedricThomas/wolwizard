package base

import (
	"context"
	"fmt"

	"github.com/CedricThomas/console/internal/boundary/out/async"
	"github.com/CedricThomas/console/internal/boundary/out/keystore"
	"github.com/CedricThomas/console/internal/controller"
	"github.com/CedricThomas/console/internal/domain"
	domainasync "github.com/CedricThomas/console/internal/domain/async"
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

func (w web) BootSelectedOS(ctx context.Context, osName domain.OSName) error {
	// Publish boot command to Redis pubsub channel
	bootCmd := domainasync.BootMessage{OSName: osName}
	if err := w.publisher.Publish(ctx, domainasync.BootChannel, bootCmd); err != nil {
		return fmt.Errorf("publish boot command: %w", err)
	}

	return nil
}
