package interactor

import (
	"auth-project/src/domain/model"
	"auth-project/src/infrastructure/authentication"
	"auth-project/src/usecase/presenter"
	"auth-project/src/usecase/repository"
	"auth-project/tools"
	"context"
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/lindell/go-burner-email-providers/burner"
)

type userInteractor struct {
	AuthRepository  repository.AuthRepository
	UserRepository  repository.UserRepository
	TokenRepository repository.TokenRepository

	UserPresenter presenter.UserPresenter

	jwtConfigurator *authentication.JwtConfigurator
}

type UserInteractor interface {
	SignUp(ctx context.Context, signUpReq *model.SignUpReq) error
	SignUpSendOTP(ctx context.Context, signUpSend2faCodeReq *model.SignUpSend2faCodeReq) error

	VerifyResetUserPasswordCode(ctx context.Context, userResetPasswordReq *model.VerifyResetUserPassword2faСodeReq) error
	SendCodeForResetUserPassword(ctx context.Context, send2faCodeForResetUserPasswordReq *model.Send2faCodeForResetUserPasswordReq) (map[string]string, error)

	GetMyProfileByID(ctx context.Context, userID string) (*model.UserGetMyProfileResp, error)

	ChangeMyPassword(ctx context.Context, reqData *model.UserChangePasswordReq, usrID string) error
	UpdateUserInfoByID(ctx context.Context, updReq *model.UserUpdateInfoData, userID string) (*model.UserUpdResp, error)
	UpdateMyselfPhone(ctx context.Context, reqData *model.UserPhoneUpdateReq, userID string) (*model.UserUpdResp, error)
	UpdateMyselfEmail(ctx context.Context, reqData *model.UserEmailUpdateReq, userID string) (*model.UserUpdResp, error)

	SignOut(ctx context.Context, sessionID string) error
	SignOutAll(ctx context.Context, usrID string) error
}

func NewUserInteractor(
	ar repository.AuthRepository, ur repository.UserRepository, tr repository.TokenRepository, p presenter.UserPresenter, jc *authentication.JwtConfigurator) UserInteractor {
	return &userInteractor{ar, ur, tr, p, jc}
}

func (ui *userInteractor) SignUp(ctx context.Context, signUpReq *model.SignUpReq) error {
	var err error

	signUpReq.Password, err = tools.VerifyPassword(signUpReq.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var verifyCodeDate *model.VerifyCodeData

	switch signUpReq.LoginType {
	case model.TokenTypePhone:
		signUpReq.Login, err = tools.VerifyPhone(signUpReq.Login)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		verifyCodeDate = &model.VerifyCodeData{
			Target:   signUpReq.Login,
			Code:     signUpReq.Code2fa,
			CodeType: model.TokenTypePhone,
			Reason:   model.TokenReasonSignUp,
		}

		_, err = ui.TokenRepository.Validate2faCode(ctx, verifyCodeDate)
		if err != nil {
			return err
		}

	case model.TokenTypeEmail:
		verifyCodeDate = &model.VerifyCodeData{
			Target:   signUpReq.Login,
			Code:     signUpReq.Code2fa,
			CodeType: model.TokenTypeEmail,
			Reason:   model.TokenReasonSignUp,
		}

		_, err = ui.TokenRepository.Validate2faCode(ctx, verifyCodeDate)
		if err != nil {
			return err
		}

	default:
		return errors.New("invalid login type")
	}

	err = ui.UserRepository.SignUpActivateUser(ctx, &model.SignUpActivateUserData{
		Login:     signUpReq.Login,
		LoginType: signUpReq.LoginType,
		Password:  signUpReq.Password,
		Referral:  signUpReq.Referral,
	})
	if err != nil {
		return err
	}

	if verifyCodeDate != nil {
		err = ui.TokenRepository.TokenSetUsed(ctx, verifyCodeDate)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ui *userInteractor) SignUpSendOTP(ctx context.Context, signUpSendOTPReq *model.SignUpSend2faCodeReq) error {
	var err error

	switch signUpSendOTPReq.LoginType {
	case model.TokenTypePhone:
		signUpSendOTPReq.Login, err = tools.VerifyPhone(signUpSendOTPReq.Login)
		if err != nil {
			return err
		}

		err = ui.UserRepository.CreatUnActivateUserIfNotExistByLogin(ctx, signUpSendOTPReq)
		if err != nil {
			return err
		}

		_, err = ui.TokenRepository.Send2faCode(ctx, &model.Send2faCodeData{
			Target:      signUpSendOTPReq.Login,
			Code2faType: model.TokenTypePhone,
			Reason:      model.TokenReasonSignUp,
		})
		if err != nil {
			return err
		}

	case model.TokenTypeEmail:
		isBurnerEmail := burner.IsBurnerEmail(signUpSendOTPReq.Login)
		if isBurnerEmail {
			return fiber.NewError(fiber.StatusInternalServerError, "invalid email")
		}

		err = ui.UserRepository.CreatUnActivateUserIfNotExistByLogin(ctx, signUpSendOTPReq)
		if err != nil {
			return err
		}

		_, err = ui.TokenRepository.Send2faCode(ctx, &model.Send2faCodeData{
			Target:      signUpSendOTPReq.Login,
			Code2faType: model.TokenTypeEmail,
			Reason:      model.TokenReasonSignUp,
		})
		if err != nil {
			return err
		}
	default:
		return errors.New("invalid login type")
	}

	return nil
}

func (ui *userInteractor) VerifyResetUserPasswordCode(ctx context.Context,
	userResetPasswordReq *model.VerifyResetUserPassword2faСodeReq) error {
	var err error

	verifyCodeDate := &model.VerifyCodeData{
		Target:   userResetPasswordReq.Target,
		Code:     userResetPasswordReq.Code2fa,
		CodeType: userResetPasswordReq.Code2faType,
		Reason:   model.TokenReasonResetPassword,
	}

	token, err := ui.TokenRepository.Validate2faCode(ctx, verifyCodeDate)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid code")
	}

	user, err := ui.UserRepository.GetUserByEmailOrPhone(ctx, token.Target)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "invalid user login")
	}

	err = ui.UserRepository.ResetUserPassword(ctx, userResetPasswordReq.NewPassword, user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func (ui *userInteractor) SendCodeForResetUserPassword(ctx context.Context,
	send2faCodeForResetUserPasswordReq *model.Send2faCodeForResetUserPasswordReq) (map[string]string, error) {

	var err error
	switch send2faCodeForResetUserPasswordReq.TargetType {
	case model.TokenTypePhone:

		send2faCodeForResetUserPasswordReq.Target, err = tools.VerifyPhone(send2faCodeForResetUserPasswordReq.Target)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		user, err := ui.UserRepository.GetUserByEmailOrPhone(ctx, send2faCodeForResetUserPasswordReq.Target)
		if err != nil {
			if err == sql.ErrNoRows {
				return map[string]string{
					"message": "OK",
				}, nil
			}
			return nil, err
		}

		_, err = ui.TokenRepository.Send2faCode(ctx, &model.Send2faCodeData{
			UserID:      user.ID,
			Target:      user.Phone,
			Code2faType: model.TokenTypePhone,
			Reason:      model.TokenReasonResetPassword,
		})
		if err != nil {
			return nil, err
		}

		return map[string]string{
			"message": "OK",
		}, nil
	case model.TokenTypeEmail:
		isBurnerEmail := burner.IsBurnerEmail(send2faCodeForResetUserPasswordReq.Target)
		if isBurnerEmail {
			return nil, fiber.NewError(fiber.StatusBadRequest, "invalid email")
		}

		user, err := ui.UserRepository.GetUserByEmailOrPhone(ctx, send2faCodeForResetUserPasswordReq.Target)
		if err != nil {
			if err == sql.ErrNoRows {
				return map[string]string{
					"message": "OK",
				}, nil
			}
			return nil, err
		}

		_, err = ui.TokenRepository.Send2faCode(ctx, &model.Send2faCodeData{
			UserID:      user.ID,
			Target:      user.Email,
			Code2faType: model.TokenTypeEmail,
			Reason:      model.TokenReasonResetPassword,
		})
		if err != nil {
			return nil, err
		}

		return map[string]string{
			"message": "OK",
		}, nil
	default:
		return nil, errors.New("target type invalid")
	}
}

func (ui *userInteractor) GetMyProfileByID(ctx context.Context, userID string) (*model.UserGetMyProfileResp, error) {
	usr, err := ui.UserRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return ui.UserPresenter.GetMyProfileByIDResp(usr), nil
}

func (ui *userInteractor) ChangeMyPassword(ctx context.Context, reqData *model.UserChangePasswordReq, usrID string) error {
	var err error
	reqData.NewPassword, err = tools.VerifyPassword(reqData.NewPassword)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, err := ui.UserRepository.GetUserByID(ctx, usrID)
	if err != nil {
		return err
	}

	if (user.IsEmailVerified || user.IsGoogleVerified || user.IsPhoneVerified) &&
		(reqData.Code2fa == "") {
		return fiber.NewError(fiber.StatusBadRequest, "this action needs 2fa")
	}

	verifyCodeDate := &model.VerifyCodeData{
		UserID: user.ID,
		Code:   reqData.Code2fa,
		Reason: model.TokenReasonVerification,
	}

	if reqData.Code2fa != "" {
		_, err = ui.TokenRepository.Validate2faCode(ctx, verifyCodeDate)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	err = ui.UserRepository.ChangeUserPasswordByID(ctx,
		&model.UserChangePasswordData{
			NewPassword: reqData.NewPassword,
			OldPassword: reqData.OldPassword,
		}, usrID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if reqData.Code2fa != "" {
		err = ui.TokenRepository.TokenSetUsed(ctx, verifyCodeDate)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return nil
}

func (ui *userInteractor) UpdateUserInfoByID(ctx context.Context, updReq *model.UserUpdateInfoData,
	userID string) (*model.UserUpdResp, error) {

	user, err := ui.UserRepository.UpdateUserInfoByID(ctx, updReq, userID)
	if err != nil {
		return nil, err
	}
	return ui.UserPresenter.UpdateUserByIDResp(user), nil
}

func (ui *userInteractor) UpdateMyselfPhone(ctx context.Context, reqData *model.UserPhoneUpdateReq,
	userID string) (*model.UserUpdResp, error) {

	var err error

	reqData.Phone, err = tools.VerifyPhone(reqData.Phone)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	verifyCodeDate := &model.VerifyCodeData{
		UserID:   userID,
		Target:   reqData.Phone,
		Code:     reqData.Code2fa,
		CodeType: model.TokenTypePhone,
		Reason:   model.TokenReasonVerification,
	}

	_, err = ui.TokenRepository.Validate2faCode(ctx, verifyCodeDate)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	user, err := ui.UserRepository.UpdateUserPhoneByID(ctx, reqData.Phone, userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err = ui.TokenRepository.TokenSetUsed(ctx, verifyCodeDate)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ui.UserPresenter.UpdateUserByIDResp(user), nil
}

func (ui *userInteractor) UpdateMyselfEmail(ctx context.Context,
	reqData *model.UserEmailUpdateReq, usrID string) (*model.UserUpdResp, error) {

	isBurnerEmail := burner.IsBurnerEmail(reqData.Email)
	if isBurnerEmail {
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid email")
	}

	verifyCodeDate := &model.VerifyCodeData{
		UserID:   usrID,
		Target:   reqData.Email,
		Code:     reqData.Code2fa,
		CodeType: model.TokenTypeEmail,
		Reason:   model.TokenReasonVerification,
	}

	_, err := ui.TokenRepository.Validate2faCode(ctx, verifyCodeDate)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	user, err := ui.UserRepository.UpdateUserEmailByID(ctx, reqData.Email, usrID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err = ui.TokenRepository.TokenSetUsed(ctx, verifyCodeDate)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ui.UserPresenter.UpdateUserByIDResp(user), nil
}

func (ui *userInteractor) SignOut(ctx context.Context, sessionID string) error {

	err := ui.UserRepository.SignOut(ctx, sessionID)
	if err != nil {
		return err
	}
	return nil
}

func (ui *userInteractor) SignOutAll(ctx context.Context, usrID string) error {

	err := ui.UserRepository.SignOutAll(ctx, usrID)
	if err != nil {
		return err
	}
	return nil
}
