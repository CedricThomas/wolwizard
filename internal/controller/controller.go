package controller

import (
	"context"

	"github.com/CedricThomas/console/internal/domain"
)

type Auth interface {
	CreateAccount(ctx context.Context, username, password string) error
	CheckAuth(ctx context.Context, username, password string) (bool, error)
	DeleteAccount(ctx context.Context, username string) error
	GenerateToken(ctx context.Context, username string) (string, error)
	ValidateToken(ctx context.Context, token string) (string, error)
	RevokeToken(ctx context.Context, token string) error
	RevokeAllTokens(ctx context.Context, username string) error
}

type Web interface {
	Auth
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
