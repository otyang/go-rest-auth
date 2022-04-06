package repository

import (
	"context"
	"github.com/google/uuid"

	"auth-project/src/domain/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, usrCrtReq *model.UserCreateReq) (string, error)
	ChangeUserPassword(ctx context.Context, reqData *model.UserChangePasswordReq) error
	UpdateUserByID(ctx context.Context, updReq *model.UserUpdateReq, userID uuid.UUID) (*model.User, error)
	UpdatePinById(ctx context.Context, updReq *model.UserPinUpdateReq, userID uuid.UUID) error
	DeletePinByID(ctx context.Context, updReq *model.UserPinDeleteReq, userID uuid.UUID) error
	SignOut(ctx context.Context, atID string) error
	SignOutAll(ctx context.Context, usrID string) error
}
