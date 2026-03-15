package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/CedricThomas/console/internal/service/token"
	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
	secretKey string
	expiry    time.Duration
}

func New(secretKey string, expirySeconds int) token.Service {
	return &jwtService{
		secretKey: secretKey,
		expiry:    time.Duration(expirySeconds) * time.Second,
	}
}

func (j *jwtService) Sign(ctx context.Context, subject string) (string, time.Duration, error) {
	claims := jwt.MapClaims{
		"sub": subject,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(j.expiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", 0, fmt.Errorf("sign token: %w", err)
	}

	return signedToken, j.expiry, nil
}

func (j *jwtService) Verify(ctx context.Context, tokenStr string) (string, bool, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return "", false, err
	}

	if token.Valid {
		subject, ok := claims["sub"].(string)
		if ok {
			return subject, true, nil
		}
	}

	return "", false, nil
}
