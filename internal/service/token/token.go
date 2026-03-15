package token

import "context"

type Service interface {
	Sign(ctx context.Context, subject string) (string, error)
	Verify(ctx context.Context, token string) (subject string, valid bool, err error)
}
