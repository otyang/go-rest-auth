package presenter

import (
	"auth-project/src/domain/model"
)

type UserPresenter interface {
	GetMyProfileByIDResp(usr *model.User) *model.UserGetMyProfileResp
	GetUserByLoginAndPasswordResp(usr *model.User) *model.User

	UpdateUserByIDResp(user *model.User) *model.UserUpdResp
	ChangeUserPasswordResp(reqData *model.UserChangePasswordReq) error
}
