package repository

import (
	"context"

	"auth-project/src/domain/model"
)

type AuthRepository interface {
	StoreTokenPair(ctx context.Context, td *model.TokenDetails, sessionID string) error
	StoreAccessToken(ctx context.Context, atd *model.AccessTokenDetails, sessionID string) error

	FetchAuth(ctx context.Context, sessionID string) (string, error)
	ValidateAccessToken(ctx context.Context, atID string, sessionID string) error
}
