package base

import (
	"context"

	"github.com/CedricThomas/console/internal/boundary/out/command"
	"github.com/CedricThomas/console/internal/controller"
)

type pcAgent struct {
	executor command.CommandExecutor
}

func NewPCAgentController(executor command.CommandExecutor) controller.PCAgent {
	return &pcAgent{
		executor: executor,
	}
}

func (pa *pcAgent) ShutdownCurrentHost(ctx context.Context) error {
	return pa.executor.Shutdown(ctx)
}
