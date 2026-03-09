package base

import (
	"context"

	"github.com/CedricThomas/console/internal/controller"
)

type pcAgent struct{}

func NewPCAgentController() controller.PCAgent {
	return pcAgent{}
}

func (pa pcAgent) ShutdownCurrentHost(_ context.Context) error {
	return nil
}
