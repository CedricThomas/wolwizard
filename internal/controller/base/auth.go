package base

import (
	"context"

	"github.com/CedricThomas/console/internal/controller"
	"github.com/CedricThomas/console/internal/service/keystore"
	"github.com/CedricThomas/console/internal/service/token"
	usecase "github.com/CedricThomas/console/internal/usecase/auth"
	basauth "github.com/CedricThomas/console/internal/usecase/auth/base"
)

type auth struct {
	authUsecase usecase.Auth
}

func NewAuthController(keystore keystore.Keystore, tokenSrv token.Service) controller.Auth {
	return &auth{
		authUsecase: basauth.New(keystore, tokenSrv),
	}
}

func newAuthController(keystore keystore.Keystore, tokenSrv token.Service) auth {
	return auth{
		authUsecase: basauth.New(keystore, tokenSrv),
	}
}

func (a *auth) CreateAccount(ctx context.Context, username, password string) error {
	return a.authUsecase.CreateAccount(ctx, username, password)
}

func (a *auth) CheckAuth(ctx context.Context, username, password string) (bool, error) {
	return a.authUsecase.CheckAuth(ctx, username, password)
}

func (a *auth) DeleteAccount(ctx context.Context, username string) error {
	return a.authUsecase.DeleteAccount(ctx, username)
}

func (a *auth) GenerateToken(ctx context.Context, username string) (string, error) {
	return a.authUsecase.GenerateToken(ctx, username)
}

func (a *auth) ValidateToken(ctx context.Context, token string) (string, error) {
	return a.authUsecase.ValidateToken(ctx, token)
}

func (a *auth) RevokeToken(ctx context.Context, token string) error {
	return a.authUsecase.RevokeToken(ctx, token)
}

func (a *auth) RevokeAllTokens(ctx context.Context, username string) error {
	return a.authUsecase.RevokeAllTokens(ctx, username)
}
