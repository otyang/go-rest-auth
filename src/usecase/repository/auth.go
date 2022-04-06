package repository

import (
	"context"

	"auth-project/src/domain/model"
)

type AuthRepository interface {
	GetUserByEmailAndPassword(ctx context.Context, email, psw string) (*model.User, error)
	StoreAuth(ctx context.Context, usrInfo *model.UserRedisSessionData, td *model.TokenDetails) error
	GenerateQrCode(ctx context.Context) ([]byte, string, error)
	FetchAuth(ctx context.Context, rtId string) (string, error)
	ValidateAccessTokenID(ctx context.Context, usrID string, atID string) error
}
