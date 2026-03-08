package base

import (
	"context"

	"github.com/CedricThomas/console/internal/controller"
)

type pcAgent struct {
}

func NewPCAgentController() controller.Web {
	return web{}
}

func (pa pcAgent) ShutdowncurrentHost(_ context.Context) error {
	return nil
}
