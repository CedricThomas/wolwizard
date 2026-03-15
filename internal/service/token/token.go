package token

import (
	"context"
	"time"
)

type Service interface {
	Sign(ctx context.Context, subject string) (string, time.Duration, error)
	Verify(ctx context.Context, token string) (subject string, valid bool, err error)
}
