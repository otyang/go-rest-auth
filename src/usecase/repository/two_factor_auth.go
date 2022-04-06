package repository

import (
	"auth-project/src/domain/model"
	"context"
)

type TwoFactorAuthRepository interface {
	SetUpTwoFactorAuthByUserID(ctx context.Context, twoFactorAuthSetUpReq *model.TwoFactorAuthSetUpReq, usrID string) error
	DeleteTwoFactorAuthByUserID(ctx context.Context, twoFactorAuthType, usrID string) error
}
