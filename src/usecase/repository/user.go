package repository

import (
	"auth-project/src/domain/model"
	"context"
)

type UserRepository interface {
	IsExitsUserByEmail(ctx context.Context, email string) (bool, error)
	IsExitsUserByIDAndPassword(ctx context.Context, userID, password string) (bool, error)

	GetUserByLoginAndPassword(ctx context.Context, login, psw string) (*model.User, error)
	GetUserByEmailOrPhone(ctx context.Context, emailOrPhone string) (*model.User, error)
	GetUserByID(ctx context.Context, usrID string) (*model.User, error)

	CreatUnActivateUserIfNotExistByLogin(ctx context.Context, data *model.SignUpSend2faCodeReq) error
	SignUpActivateUser(ctx context.Context, signUpReq *model.SignUpActivateUserData) error

	ResetUserPassword(ctx context.Context, newPassword, usrID string) error

	ChangeUserPasswordByID(ctx context.Context, data *model.UserChangePasswordData, usrID string) error

	UpdateUserInfoByID(ctx context.Context, updReq *model.UserUpdateInfoData, userID string) (*model.User, error)
	UpdateUserEmailByID(ctx context.Context, email string, userID string) (*model.User, error)
	UpdateUserPhoneByID(ctx context.Context, phone string, userID string) (*model.User, error)

	SignOut(ctx context.Context, sessionID string) error
	SignOutAll(ctx context.Context, usrID string) error
}
