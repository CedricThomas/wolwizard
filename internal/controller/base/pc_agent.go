package base

import (
	"context"

	"github.com/CedricThomas/console/internal/controller"
)

type pcAgent struct{}

func NewPCAgentController() controller.PCAgent {
	return pcAgent{}
}

func (pa pcAgent) ShutdowncurrentHost(_ context.Context) error {
	return nil
}
