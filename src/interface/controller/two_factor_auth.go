package controller

import (
	"auth-project/src/domain/model"
	"auth-project/src/usecase/interactor"
	"github.com/gofiber/fiber/v2"
)

type twoFactorAuthController struct {
	twoFactorAuthInteractor interactor.TwoFactorAuthInteractor
}

type TwoFactorAuthController interface {
	ReSendTwoFactorAuthCode(ctx *fiber.Ctx) error
	VerifyTwoFactorAuthCode(ctx *fiber.Ctx) error

	GenerateGoogleTwoFactorAuthQrCode(ctx *fiber.Ctx) error
	SetUpTwoFactorAuth(ctx *fiber.Ctx) error

	DeleteTwoFactorAuth(ctx *fiber.Ctx) error
}

func NewTwoFactorAuthController(ti interactor.TwoFactorAuthInteractor) TwoFactorAuthController {
	return &twoFactorAuthController{ti}
}

// ReSendTwoFactorAuthCode checks the user's entity and re-sends the 2fa code depending on it
func (tc *twoFactorAuthController) ReSendTwoFactorAuthCode(ctx *fiber.Ctx) error {

	usrID, ok := ctx.Context().Value("token_user_id").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "context value type invalid")
	}

	resp, err := tc.twoFactorAuthInteractor.ReSendTwoFactorAuthCode(ctx.Context(), usrID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// VerifyTwoFactorAuthCode accepts the 2fa code, verify him and returns a token to authorize a client
func (tc *twoFactorAuthController) VerifyTwoFactorAuthCode(ctx *fiber.Ctx) error {

	var verify2faCodeReq *model.Verify2faCodeReq
	err := ctx.BodyParser(&verify2faCodeReq)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	usrID, ok := ctx.Context().Value("token_user_id").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "context value type invalid")
	}

	usrInfo := &model.UserSessionData{
		UserID:    usrID,
		UserAgent: string(ctx.Request().Header.UserAgent()),
		ClientIp:  ctx.Context().RemoteAddr().String(),
	}

	details, err := tc.twoFactorAuthInteractor.VerifyTwoFactorAuthCode(ctx.Context(), verify2faCodeReq, usrInfo)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]string{
		"access_token":  details.AccessToken,
		"refresh_token": details.RefreshToken,
	})
}

// GenerateGoogleTwoFactorAuthQrCode generates qr code for google 2fa
func (tc *twoFactorAuthController) GenerateGoogleTwoFactorAuthQrCode(ctx *fiber.Ctx) error {

	usrID, ok := ctx.Context().Value("token_user_id").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "context value type invalid")
	}

	resp, err := tc.twoFactorAuthInteractor.GenerateGoogleTwoFactorAuthQrCode(ctx.Context(), usrID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// SetUpTwoFactorAuth accepts 2fa code, verify it and sets it to the user
func (tc *twoFactorAuthController) SetUpTwoFactorAuth(ctx *fiber.Ctx) error {

	var twoFactorAuthSetUpReq model.TwoFactorAuthSetUpReq
	err := ctx.BodyParser(&twoFactorAuthSetUpReq)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	usrID, ok := ctx.Context().Value("token_user_id").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "context value type invalid")
	}

	err = tc.twoFactorAuthInteractor.SetUpTwoFactorAuthByUserID(ctx.Context(), &twoFactorAuthSetUpReq,
		usrID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "OK",
	})
}

// DeleteTwoFactorAuth accepts user password, verify it and delete 2fa
func (tc *twoFactorAuthController) DeleteTwoFactorAuth(ctx *fiber.Ctx) error {

	var twoFactorAuthDeleteReq model.TwoFactorAuthDeleteReq
	err := ctx.BodyParser(&twoFactorAuthDeleteReq)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	usrID, ok := ctx.Context().Value("token_user_id").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "context value type invalid")
	}

	err = tc.twoFactorAuthInteractor.DeleteTwoFactorAuthByUserID(ctx.Context(), &twoFactorAuthDeleteReq,
		usrID)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "OK",
	})
}
