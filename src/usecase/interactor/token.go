package interactor

import (
	"auth-project/src/domain/model"
	"auth-project/src/usecase/presenter"
	"auth-project/src/usecase/repository"
	"auth-project/tools"
	"context"
	"github.com/gofiber/fiber/v2"
)

type tokenInteractor struct {
	TokenRepository repository.TokenRepository
	UserRepository  repository.UserRepository

	TokenPresenter presenter.TokenPresenter
}

type TokenInteractor interface {
	Validate2faCode(ctx context.Context, verifyCodeDate *model.VerifyCodeData) (*model.Token, error)
	Send2faCode(ctx context.Context, usrID string) (map[string]interface{}, error)
	SendTarget2faCode(ctx context.Context, sendTarget2faCodeReq *model.SendTarget2faCodeReq, usrID string) (map[string]interface{}, error)
	TokenSetUsed(ctx context.Context, verifyCodeDate *model.VerifyCodeData) error
}

func NewTokenInteractor(
	tr repository.TokenRepository, ur repository.UserRepository, p presenter.TokenPresenter) TokenInteractor {
	return &tokenInteractor{tr, ur, p}
}

func (ti *tokenInteractor) Validate2faCode(ctx context.Context, verifyCodeDate *model.VerifyCodeData) (*model.Token, error) {

	token, err := ti.TokenRepository.Validate2faCode(ctx, verifyCodeDate)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (ti *tokenInteractor) Send2faCode(ctx context.Context, usrID string) (map[string]interface{}, error) {

	usr, err := ti.UserRepository.GetUserByID(ctx, usrID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	switch {
	case usr.IsGoogleVerified:
		return map[string]interface{}{
			"code_2fa_type": model.TokenTypeGoogle,
		}, nil

	case usr.IsPhoneVerified:
		target, err := ti.TokenRepository.Send2faCode(ctx, &model.Send2faCodeData{
			UserID:      usr.ID,
			Target:      usr.Phone,
			Code2faType: model.TokenTypePhone,
			Reason:      model.TokenReasonVerification,
		})
		if err != nil {
			return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return map[string]interface{}{
			"code_2fa_type":   model.TokenTypePhone,
			"code_2fa_target": target,
		}, nil

	case usr.IsEmailVerified:
		target, err := ti.TokenRepository.Send2faCode(ctx, &model.Send2faCodeData{
			UserID:      usr.ID,
			Target:      usr.Email,
			Code2faType: model.TokenTypeEmail,
			Reason:      model.TokenReasonVerification,
		})
		if err != nil {
			return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return map[string]interface{}{
			"code_2fa_type":   model.TokenTypeEmail,
			"code_2fa_target": target,
		}, nil

	default:
		return nil, fiber.NewError(fiber.StatusInternalServerError, "two-factor authentication disabled")
	}
}

func (ti *tokenInteractor) SendTarget2faCode(ctx context.Context, sendTarget2faCodeReq *model.SendTarget2faCodeReq,
	usrID string) (map[string]interface{}, error) {

	usr, err := ti.UserRepository.GetUserByID(ctx, usrID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	switch sendTarget2faCodeReq.Code2faType {
	case model.TokenTypePhone:

		sendTarget2faCodeReq.Target, err = tools.VerifyPhone(sendTarget2faCodeReq.Target)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		target, err := ti.TokenRepository.Send2faCode(ctx, &model.Send2faCodeData{
			UserID:      usr.ID,
			Target:      sendTarget2faCodeReq.Target,
			Code2faType: model.TokenTypePhone,
			Reason:      model.TokenReasonVerification,
		})
		if err != nil {
			return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return map[string]interface{}{
			"code_2fa_type":   model.TokenTypePhone,
			"code_2fa_target": target,
		}, nil

	case model.TokenTypeEmail:
		target, err := ti.TokenRepository.Send2faCode(ctx, &model.Send2faCodeData{
			UserID:      usr.ID,
			Target:      sendTarget2faCodeReq.Target,
			Code2faType: model.TokenTypeEmail,
			Reason:      model.TokenReasonVerification,
		})
		if err != nil {
			return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return map[string]interface{}{
			"code_2fa_type":   model.TokenTypeEmail,
			"code_2fa_target": target,
		}, nil

	default:
		return nil, fiber.NewError(fiber.StatusBadRequest, "otp type invalid")
	}

}

func (ti *tokenInteractor) TokenSetUsed(ctx context.Context, verifyCodeDate *model.VerifyCodeData) error {

	err := ti.TokenRepository.TokenSetUsed(ctx, verifyCodeDate)
	if err != nil {
		return err
	}

	return nil
}
