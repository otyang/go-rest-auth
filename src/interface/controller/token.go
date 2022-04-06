package controller

import (
	"auth-project/src/domain/model"
	"auth-project/src/usecase/interactor"
	"github.com/gofiber/fiber/v2"
)

type tokenController struct {
	tokenInteractor interactor.TokenInteractor
}

type TokenController interface {
	Send2faCode(ctx *fiber.Ctx) error
	SendTarget2faCode(ctx *fiber.Ctx) error
}

func NewTokenController(ai interactor.TokenInteractor) TokenController {
	return &tokenController{ai}
}

// Send2faCode checks the user's entity and sends the code depending on it
func (tc *tokenController) Send2faCode(ctx *fiber.Ctx) error {

	usrID, ok := ctx.Context().Value("token_user_id").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "context value type invalid")
	}

	resp, err := tc.tokenInteractor.Send2faCode(ctx.Context(), usrID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// SendTarget2faCode accepts the target and type code, and sends 2fa code
func (tc *tokenController) SendTarget2faCode(ctx *fiber.Ctx) error {

	var sendTarget2faCodeReq model.SendTarget2faCodeReq
	err := ctx.BodyParser(&sendTarget2faCodeReq)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	usrID, ok := ctx.Context().Value("token_user_id").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "context value type invalid")
	}

	resp, err := tc.tokenInteractor.SendTarget2faCode(ctx.Context(), &sendTarget2faCodeReq, usrID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}
