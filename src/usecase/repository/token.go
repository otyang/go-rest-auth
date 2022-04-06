package repository

import (
	"auth-project/src/domain/model"
	"context"
)

type TokenRepository interface {
	Validate2faCode(ctx context.Context, VerifyCodeDate *model.VerifyCodeData) (*model.Token, error)
	Send2faCode(ctx context.Context, Send2faCodeDate *model.Send2faCodeData) (string, error)
	TokenSetUsed(ctx context.Context, verifyCodeDate *model.VerifyCodeData) error
}
