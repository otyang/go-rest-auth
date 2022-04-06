package interactor

import (
	"auth-project/src/domain/model"
	"auth-project/src/infrastructure/authentication"
	"auth-project/src/usecase/presenter"
	"auth-project/src/usecase/repository"
	"context"
	"github.com/gofiber/fiber/v2"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"time"
)

type twoFactorAuthInteractor struct {
	AuthRepository          repository.AuthRepository
	SessionRepository       repository.SessionRepository
	TwoFactorAuthRepository repository.TwoFactorAuthRepository
	UserRepository          repository.UserRepository
	TokenRepository         repository.TokenRepository

	TwoFactorAuthPresenter presenter.TwoFactorAuthPresenter

	jwtConfigurator *authentication.JwtConfigurator
}

type TwoFactorAuthInteractor interface {
	ReSendTwoFactorAuthCode(ctx context.Context, usrID string) (map[string]string, error)
	VerifyTwoFactorAuthCode(ctx context.Context, verify2faCodeReq *model.Verify2faCodeReq, usrInfo *model.UserSessionData) (*model.TokenDetails, error)
	GenerateGoogleTwoFactorAuthQrCode(ctx context.Context, usrID string) (map[string]interface{}, error)

	SetUpTwoFactorAuthByUserID(ctx context.Context, twoFactorAuthSetUpReq *model.TwoFactorAuthSetUpReq, usrID string) error
	DeleteTwoFactorAuthByUserID(ctx context.Context, twoFactorAuthDeleteReq *model.TwoFactorAuthDeleteReq, usrID string) error
}

func NewTwoFactorAuthInteractor(
	ar repository.AuthRepository, sr repository.SessionRepository, tfr repository.TwoFactorAuthRepository, ur repository.UserRepository, tr repository.TokenRepository, tp presenter.TwoFactorAuthPresenter, jc *authentication.JwtConfigurator) TwoFactorAuthInteractor {
	return &twoFactorAuthInteractor{ar, sr, tfr, ur, tr, tp, jc}
}

func (ti *twoFactorAuthInteractor) ReSendTwoFactorAuthCode(ctx context.Context, usrID string) (map[string]string, error) {

	usr, err := ti.UserRepository.GetUserByID(ctx, usrID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	switch {
	case usr.IsGoogleVerified:
		return nil, fiber.NewError(fiber.StatusConflict, "the user has two factor authentication using Google")

	case usr.IsPhoneVerified:
		target, err := ti.TokenRepository.Send2faCode(ctx, &model.Send2faCodeData{
			UserID:      usr.ID,
			Target:      usr.Phone,
			Code2faType: model.TokenTypePhone,
			Reason:      model.TokenReasonTwoFactorAuth,
		})
		if err != nil {
			return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return map[string]string{
			"2fa_type":   model.TokenTypePhone,
			"2fa_target": target,
		}, nil

	case usr.IsEmailVerified:
		target, err := ti.TokenRepository.Send2faCode(ctx, &model.Send2faCodeData{
			UserID:      usr.ID,
			Target:      usr.Email,
			Code2faType: model.TokenTypeEmail,
			Reason:      model.TokenReasonTwoFactorAuth,
		})
		if err != nil {
			return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return map[string]string{
			"2fa_type":   model.TokenTypeEmail,
			"2fa_target": target,
		}, nil

	default:
		return nil, fiber.NewError(fiber.StatusInternalServerError, "two factor auth deactivate")
	}
}

func (ti *twoFactorAuthInteractor) VerifyTwoFactorAuthCode(ctx context.Context, verify2faCodeReq *model.Verify2faCodeReq,
	usrInfo *model.UserSessionData) (*model.TokenDetails, error) {
	var err error

	switch verify2faCodeReq.Code2faType {
	case model.TokenTypeGoogle:
		user, err := ti.UserRepository.GetUserByID(ctx, usrInfo.UserID)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusInternalServerError, "invalid token")
		}
		err = authentication.VerifyGoogleTwoFactorAuthCode(verify2faCodeReq.Code2fa, user.GoogleSecret)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	case model.TokenTypePhone, model.TokenTypeEmail:
		_, err = ti.TokenRepository.Validate2faCode(ctx, &model.VerifyCodeData{
			UserID: usrInfo.UserID,
			Code:   verify2faCodeReq.Code2fa,
			Reason: model.TokenReasonTwoFactorAuth,
		})
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	default:
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid token type")
	}
	sessionID, err := gonanoid.New()
	if err != nil {
		return nil, err
	}

	details, err := ti.jwtConfigurator.GenerateTokenPair(usrInfo.UserID, sessionID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err = ti.AuthRepository.StoreTokenPair(ctx, details, sessionID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	ses := &model.Session{
		SessionID: sessionID,
		UserAgent: usrInfo.UserAgent,
		ClientIP:  usrInfo.ClientIp,
		ExpiresAT: time.Unix(details.RtExpires, 0).UTC(),
		UserID:    usrInfo.UserID,
	}

	err = ti.SessionRepository.InsertSession(context.Background(), ses)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if verify2faCodeReq.Code2faType == model.TokenTypePhone || verify2faCodeReq.Code2faType == model.TokenTypeEmail {
		err = ti.TokenRepository.TokenSetUsed(ctx, &model.VerifyCodeData{
			UserID: usrInfo.UserID,
			Code:   verify2faCodeReq.Code2fa,
			Reason: model.TokenReasonTwoFactorAuth,
		})
		if err != nil {
			return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return details, nil
}

func (ti *twoFactorAuthInteractor) GenerateGoogleTwoFactorAuthQrCode(ctx context.Context, usrID string) (map[string]interface{}, error) {

	user, err := ti.UserRepository.GetUserByID(ctx, usrID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if user.GoogleSecret != "" {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "the new qr code can be connected if the old one is disabled")
	}

	qrCodeByte, secret, err := authentication.GenerateGoogleTwoFactorAuthQrCode(user.Email)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return map[string]interface{}{
		"qr_code": qrCodeByte,
		"secret":  secret,
	}, nil
}

func (ti *twoFactorAuthInteractor) SetUpTwoFactorAuthByUserID(ctx context.Context,
	twoFactorAuthSetUpReq *model.TwoFactorAuthSetUpReq, usrID string) error {
	var err error
	var verifyCodeData *model.VerifyCodeData

	switch twoFactorAuthSetUpReq.Code2faType {
	case model.TokenTypeGoogle:
		err = authentication.VerifyGoogleTwoFactorAuthCode(twoFactorAuthSetUpReq.Code2fa, twoFactorAuthSetUpReq.Secret)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid token")
		}

	case model.TokenTypePhone:
		user, err := ti.UserRepository.GetUserByID(ctx, usrID)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if user.Phone == ""{
			return fiber.NewError(fiber.StatusBadRequest, "user phone missing")
		}

		verifyCodeData = &model.VerifyCodeData{
			UserID:   user.ID,
			Target:   user.Phone,
			Code:     twoFactorAuthSetUpReq.Code2fa,
			CodeType: model.TokenTypePhone,
			Reason:   model.TokenReasonVerification,
		}

		_, err = ti.TokenRepository.Validate2faCode(ctx, verifyCodeData)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

	case model.TokenTypeEmail:
		user, err := ti.UserRepository.GetUserByID(ctx, usrID)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if user.Email == ""{
			return fiber.NewError(fiber.StatusBadRequest, "user email missing")
		}

		verifyCodeData = &model.VerifyCodeData{
			UserID:   user.ID,
			Target:   user.Email,
			Code:     twoFactorAuthSetUpReq.Code2fa,
			CodeType: model.TokenTypeEmail,
			Reason:   model.TokenReasonVerification,
		}

		_, err = ti.TokenRepository.Validate2faCode(ctx, verifyCodeData)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

	default:
		return fiber.NewError(fiber.StatusInternalServerError, "invalid two factor auth type")
	}

	err = ti.TwoFactorAuthRepository.SetUpTwoFactorAuthByUserID(ctx, twoFactorAuthSetUpReq, usrID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if verifyCodeData != nil {
		err = ti.TokenRepository.TokenSetUsed(ctx, verifyCodeData)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return nil
}

func (ti *twoFactorAuthInteractor) DeleteTwoFactorAuthByUserID(ctx context.Context,
	twoFactorAuthDeleteReq *model.TwoFactorAuthDeleteReq, usrID string) error {

	_, err := ti.UserRepository.IsExitsUserByIDAndPassword(ctx, usrID,
		twoFactorAuthDeleteReq.Password)
	if err != nil {
		return err
	}

	err = ti.TwoFactorAuthRepository.DeleteTwoFactorAuthByUserID(ctx, twoFactorAuthDeleteReq.Type, usrID)
	if err != nil {
		return err
	}
	return nil
}
