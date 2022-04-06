package controller

import (
	"auth-project/src/domain/model"
	"auth-project/src/usecase/interactor"
	"auth-project/tools"
	"github.com/gofiber/fiber/v2"
)

type authController struct {
	authInteractor interactor.AuthInteractor
}

type AuthController interface {
	Authenticate(ctx *fiber.Ctx) error
	RefreshToken(ctx *fiber.Ctx) error

	ValidateAccessToken(ctx *fiber.Ctx) error
	ValidateTwoFactorAuthToken(ctx *fiber.Ctx) error
}

func NewAuthController(ai interactor.AuthInteractor) AuthController {
	return &authController{ai}
}

// Authenticate accepts the user form data and returns a token to authorize a client
func (ac *authController) Authenticate(ctx *fiber.Ctx) error {

	var authReq model.AuthenticationReq
	err := ctx.BodyParser(&authReq)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	usrInfo := &model.UserSessionData{
		UserAgent: string(ctx.Request().Header.UserAgent()),
		ClientIp:  ctx.Context().RemoteAddr().String(),
	}

	resp, err := ac.authInteractor.Authenticate(ctx.Context(), &authReq, usrInfo)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// RefreshToken accepts the refresh token, verify and returns a token to authorize a client
func (ac *authController) RefreshToken(ctx *fiber.Ctx) error {

	bearerToken, err := tools.ParseAndCheckToken(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	usrInfo := &model.UserSessionData{
		UserAgent: string(ctx.Request().Header.UserAgent()),
		ClientIp:  ctx.Context().RemoteAddr().String(),
	}

	resp, err := ac.authInteractor.RefreshToken(ctx.Context(), usrInfo, bearerToken)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// ValidateAccessToken gets the access token and verify him
func (ac *authController) ValidateAccessToken(ctx *fiber.Ctx) error {

	bearerToken, err := tools.ParseAndCheckToken(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	claims, err := ac.authInteractor.ValidateAccessToken(ctx.Context(), bearerToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	// set user id from token in context
	if claims.UserID != "" && claims.AtID != "" {
		ctx.Context().SetUserValue("token_user_id", claims.UserID)
		ctx.Context().SetUserValue("token_session_id", claims.SessionID)
	} else {
		return fiber.NewError(fiber.StatusBadRequest, "claims is missing")
	}

	return nil
}

// ValidateTwoFactorAuthToken gets the access 2fa token and verify him
func (ac *authController) ValidateTwoFactorAuthToken(ctx *fiber.Ctx) error {

	bearerToken, err := tools.ParseAndCheckToken(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	claims, err := ac.authInteractor.ValidateTwoFactorAuthToken(ctx.Context(), bearerToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	// set user id from token in context
	if claims.UserID != "" && claims.AtID != "" {
		ctx.Context().SetUserValue("token_user_id", claims.UserID)
	} else {
		return fiber.NewError(fiber.StatusBadRequest, "claims is missing")
	}

	return nil
}
