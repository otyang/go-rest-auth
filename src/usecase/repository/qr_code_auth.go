package repository

import (
	"auth-project/src/domain/model"
	"context"
)

type QrCodeAuthRepository interface {
	GenerateQrCode(ctx context.Context) ([]byte, string, error)
	CreateAuthToken(ctx context.Context, usrID string) (string, error)
	GetUserByQrCodeAuthToken(ctx context.Context, token string) (*model.User, error)
}
