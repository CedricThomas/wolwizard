package token

//go:generate mockgen -source=token.go -destination=mock/token.go -package=mock -mock_names=Service=MockService
import (
	"context"
	"time"
)

type Service interface {
	Sign(ctx context.Context, subject string) (string, time.Duration, error)
	Verify(ctx context.Context, token string) (subject string, valid bool, err error)
}
