package presenter

import (
	"auth-project/src/domain/model"
)

type AuthPresenter interface {
	GetUserByEmailAndPasswordResp(usr *model.User) *model.User
}
