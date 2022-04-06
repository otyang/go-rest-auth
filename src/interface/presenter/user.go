package presenter

import (
	"auth-project/src/domain/model"
)

type userPresenter struct {
}

type UserPresenter interface {
	GetMyProfileByIDResp(usr *model.User) *model.UserGetMyProfileResp
	GetUserByLoginAndPasswordResp(usr *model.User) *model.User

	UpdateUserByIDResp(user *model.User) *model.UserUpdResp
	ChangeUserPasswordResp(reqData *model.UserChangePasswordReq) error
}

func NewUserPresenter() UserPresenter {
	return &userPresenter{}
}

func (up *userPresenter) GetMyProfileByIDResp(usr *model.User) *model.UserGetMyProfileResp {
	resp := &model.UserGetMyProfileResp{
		ID:               usr.ID,
		FullName:         usr.FullName,
		UserName:         usr.UserName,
		Email:            usr.Email,
		Phone:            usr.Phone,
		ReferralLink:     usr.ReferralLink,
		Role:             usr.Role,
		IsPhoneVerified:  usr.IsPhoneVerified,
		IsEmailVerified:  usr.IsEmailVerified,
		IsGoogleVerified: usr.IsGoogleVerified,
		Referral:         usr.Referral,
		CreatedAt:        usr.CreatedAt,
	}
	if usr.ReferralUser != nil {
		resp.ReferralUser = &model.UserReferralGet{
			ID: usr.ReferralUser.ID,
		}
	}
	return resp

}

func (up *userPresenter) GetUserByLoginAndPasswordResp(usr *model.User) *model.User {
	return usr
}

func (up *userPresenter) UpdateUserByIDResp(user *model.User) *model.UserUpdResp {
	return &model.UserUpdResp{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		Phone:    user.Phone,
		Referral: user.Referral,
	}
}

func (up *userPresenter) ChangeUserPasswordResp(reqData *model.UserChangePasswordReq) (err error) {
	return err
}
