package presenter

import (
	"auth-project/src/domain/model"
)

type userPresenter struct {
}

type UserPresenter interface {
	CreateUserResp(usrCrtReq *model.UserCreateReq) (string, error)
	UpdateUserByIDResp(user *model.User) *model.UserUpdResp
	ChangeUserPasswordResp(reqData *model.UserChangePasswordReq) error
}

func NewUserPresenter() UserPresenter {
	return &userPresenter{}
}

func (up *userPresenter) CreateUserResp(usrCrtReq *model.UserCreateReq) (usrID string, err error) {
	return
}

func (up *userPresenter) UpdateUserByIDResp(user *model.User) *model.UserUpdResp {
	return &model.UserUpdResp{
		ID:           user.ID,
		FullName:     user.FullName,
		Email:        user.Email,
		Phone:        user.Phone,
		ReferralLink: user.ReferralLink,
	}
}

func (up *userPresenter) ChangeUserPasswordResp(reqData *model.UserChangePasswordReq) (err error) {
	return err
}
