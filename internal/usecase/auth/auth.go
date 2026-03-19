package auth

//go:generate mockgen -source=auth.go -destination=mock/auth.go -package=mock -mock_names=Auth=MockAuth
import "context"

type Auth interface {
	CreateAccount(ctx context.Context, username, password string) error
	CheckAuth(ctx context.Context, username, password string) (bool, error)
	DeleteAccount(ctx context.Context, username string) error
	GenerateToken(ctx context.Context, username string) (string, error)
	ValidateToken(ctx context.Context, token string) (string, error)
	RevokeToken(ctx context.Context, token string) error
	RevokeAllTokens(ctx context.Context, username string) error
}
