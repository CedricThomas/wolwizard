package controller

import (
	"context"

	"github.com/CedricThomas/console/internal/domain"
)

type Web interface {
	BoostSelectedOS(ctx context.Context, osName domain.OSName) error
}

type PCAgent interface {
	ShutdowncurrentHost(ctx context.Context) error
}
