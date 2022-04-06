package presenter

import (
	"auth-project/src/domain/model"
)

type authPresenter struct {
}

type AuthPresenter interface {
	GetUserByEmailAndPasswordResp(usr *model.User) *model.User
}

func NewAuthPresenter() AuthPresenter {
	return &authPresenter{}
}

func (ap *authPresenter) GetUserByEmailAndPasswordResp(usr *model.User) *model.User {
	return usr
}
