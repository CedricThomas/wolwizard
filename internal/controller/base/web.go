package base

import (
	"context"

	"github.com/CedricThomas/console/internal/boundary/out/async"
	"github.com/CedricThomas/console/internal/boundary/out/keystore"
	"github.com/CedricThomas/console/internal/controller"
	"github.com/CedricThomas/console/internal/domain"
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

func (w web) BoostSelectedOS(_ context.Context, _ domain.OSName) error {
	return nil
}
