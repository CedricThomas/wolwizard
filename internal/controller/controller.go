package controller

import (
	"context"

	"github.com/CedricThomas/console/internal/domain"
	"github.com/CedricThomas/console/internal/domain/async"
)

type Web interface {
	BootSelectedOS(ctx context.Context, osName domain.OSName) error
}

type RaspberryAgent interface {
	ExecuteBootMessage(ctx context.Context, bootMessage async.BootMessage) error
}

type PCAgent interface {
	ShutdowncurrentHost(ctx context.Context) error
}
