package base

import (
	"context"
	"fmt"

	"github.com/CedricThomas/console/internal/service/keystore"
	"github.com/CedricThomas/console/internal/service/token"
	"github.com/CedricThomas/console/internal/usecase/auth"

	"golang.org/x/crypto/bcrypt"
)

const (
	keyPrefix        = "auth:"
	userKeySuffix    = ":user:"
	tokenKeySuffix   = ":token:"
	passwordHashCost = 10
)

type authUsecase struct {
	keystore keystore.Keystore
	tokenSrv token.Service
}

func New(ks keystore.Keystore, ts token.Service) auth.Auth {
	return &authUsecase{
		keystore: ks,
		tokenSrv: ts,
	}
}

func (a *authUsecase) CreateAccount(ctx context.Context, username, password string) error {
	if username == "" {
		return fmt.Errorf("username is required")
	}

	exists, err := a.keystore.Exists(ctx, a.userKey(username))
	if err != nil {
		return fmt.Errorf("check existing user: %w", err)
	}
	if exists {
		return fmt.Errorf("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), passwordHashCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	err = a.keystore.Set(ctx, a.userKey(username), string(hashedPassword))
	if err != nil {
		return fmt.Errorf("store user: %w", err)
	}

	return nil
}

func (a *authUsecase) CheckAuth(ctx context.Context, username, password string) (bool, error) {
	hashedPassword, err := a.keystore.Get(ctx, a.userKey(username))
	if err != nil {
		return false, fmt.Errorf("get user: %w", err)
	}

	if hashedPassword == "" {
		return false, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, fmt.Errorf("compare password: %w", err)
	}

	return true, nil
}

func (a *authUsecase) DeleteAccount(ctx context.Context, username string) error {
	err := a.keystore.Delete(ctx, a.userKey(username))
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	tokenKeys, err := a.keystore.Keys(ctx, a.tokenKeyPattern())
	if err != nil {
		return fmt.Errorf("find user tokens: %w", err)
	}

	for _, tokenKey := range tokenKeys {
		tokenValue, err := a.keystore.Get(ctx, tokenKey)
		if err != nil {
			continue
		}
		if tokenValue == username {
			a.keystore.Delete(ctx, tokenKey)
		}
	}

	return nil
}

func (a *authUsecase) GenerateToken(ctx context.Context, username string) (string, error) {
	exists, err := a.keystore.Exists(ctx, a.userKey(username))
	if err != nil {
		return "", fmt.Errorf("check user: %w", err)
	}
	if !exists {
		return "", fmt.Errorf("user not found")
	}

	tokenStr, expiry, err := a.tokenSrv.Sign(ctx, username)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	err = a.keystore.SetWithTTL(ctx, a.tokenKey(tokenStr), username, expiry)
	if err != nil {
		return "", fmt.Errorf("store token: %w", err)
	}

	return tokenStr, nil
}

func (a *authUsecase) ValidateToken(ctx context.Context, tokenStr string) (string, error) {
	username, valid, err := a.tokenSrv.Verify(ctx, tokenStr)
	if err != nil {
		return "", fmt.Errorf("verify token: %w", err)
	}

	if !valid || username == "" {
		return "", fmt.Errorf("invalid token")
	}

	exists, err := a.keystore.Exists(ctx, a.tokenKey(tokenStr))
	if err != nil {
		return "", fmt.Errorf("check token: %w", err)
	}

	if !exists {
		return "", fmt.Errorf("token revoked")
	}

	return username, nil
}

func (a *authUsecase) RevokeToken(ctx context.Context, token string) error {
	err := a.keystore.Delete(ctx, a.tokenKey(token))
	if err != nil {
		return fmt.Errorf("revoke token: %w", err)
	}

	return nil
}

func (a *authUsecase) RevokeAllTokens(ctx context.Context, username string) error {
	tokenKeys, err := a.keystore.Keys(ctx, a.tokenKeyPattern())
	if err != nil {
		return fmt.Errorf("find user tokens: %w", err)
	}

	for _, tokenKey := range tokenKeys {
		tokenValue, err := a.keystore.Get(ctx, tokenKey)
		if err != nil {
			continue
		}
		if tokenValue == username {
			err = a.keystore.Delete(ctx, tokenKey)
			if err != nil {
				return fmt.Errorf("revoke token: %w", err)
			}
		}
	}

	return nil
}

func (a *authUsecase) userKey(username string) string {
	return keyPrefix + userKeySuffix + username
}

func (a *authUsecase) tokenKey(token string) string {
	return keyPrefix + tokenKeySuffix + token
}

func (a *authUsecase) tokenKeyPattern() string {
	return keyPrefix + tokenKeySuffix + "*"
}
