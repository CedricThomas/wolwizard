package controller

import (
	"context"

	"github.com/CedricThomas/console/internal/domain"
)

type Web interface {
	SendAsyncBootCommand(ctx context.Context, osName domain.OSName) error
	SendAsyncShutdownCommand(ctx context.Context) error
	ProcessMetrics(ctx context.Context, metrics domain.Metrics) error
}

type RaspberryAgent interface {
	WakeUpPCAgent(ctx context.Context, osName domain.OSName) error
}

type PCAgent interface {
	ShutdownCurrentHost(ctx context.Context) error
	SendCurrentHostAsyncMetrics(ctx context.Context) error
}
