package presenter

import (
	"auth-project/src/domain/model"
)

type UserPresenter interface {
	CreateUserResp(usrCrtReq *model.UserCreateReq) (string, error)
	UpdateUserByIDResp(user *model.User) *model.UserUpdResp
	ChangeUserPasswordResp(reqData *model.UserChangePasswordReq) error
}
