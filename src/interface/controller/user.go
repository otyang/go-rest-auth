package controller

//
import (
	"auth-project/src/domain/model"
	"auth-project/src/usecase/interactor"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type userController struct {
	userInteractor interactor.UserInteractor
}

type UserController interface {
	CreateUser(ctx *fiber.Ctx) error
	ChangeUserPassword(ctx *fiber.Ctx) error
	UpdatePinYourself(ctx *fiber.Ctx) error
	UpdateYourself(ctx *fiber.Ctx) error
	DeletePinYourself(ctx *fiber.Ctx) error
	SignOut(ctx *fiber.Ctx) error
	SignOutAll(ctx *fiber.Ctx) error
}

func NewUserController(ui interactor.UserInteractor) UserController {
	return &userController{ui}
}

func (uc *userController) CreateUser(ctx *fiber.Ctx) error {

	var usrCrtReq *model.UserCreateReq
	err := ctx.BodyParser(&usrCrtReq)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var usrID string
	usrID, err = uc.userInteractor.CreateUser(ctx.Context(), usrCrtReq)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]string{
		"user_id": usrID,
	})
}

// ChangeUserPassword method change old password on new password
func (uc *userController) ChangeUserPassword(ctx *fiber.Ctx) error {

	var reqData *model.UserChangePasswordReq
	err := ctx.BodyParser(&reqData)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = uc.userInteractor.ChangeUserPassword(ctx.Context(), reqData)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "OK",
	})
}

func (uc *userController) UpdatePinYourself(ctx *fiber.Ctx) error {

	var reqData *model.UserPinUpdateReq
	err := ctx.BodyParser(&reqData)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	usrIDInterface := ctx.Context().Value("token_user_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	usrID, err := uuid.Parse(usrIDInterface.(string))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = uc.userInteractor.UpdatePinById(ctx.Context(), reqData, usrID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "OK",
	})
}

func (uc *userController) DeletePinYourself(ctx *fiber.Ctx) error {

	var reqData *model.UserPinDeleteReq
	err := ctx.BodyParser(&reqData)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	usrIDInterface := ctx.Context().Value("token_user_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	usrID, err := uuid.Parse(usrIDInterface.(string))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = uc.userInteractor.DeletePinByID(ctx.Context(), reqData, usrID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "OK",
	})
}

func (uc *userController) UpdateYourself(ctx *fiber.Ctx) error {

	var reqData *model.UserUpdateReq
	err := ctx.BodyParser(&reqData)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	usrIDInterface := ctx.Context().Value("token_user_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	usrID, err := uuid.Parse(usrIDInterface.(string))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	resp, err := uc.userInteractor.UpdateUserByID(ctx.Context(), reqData, usrID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

//SignOut invalidates a user from server-side using the jwt Blocklist.
func (uc *userController) SignOut(ctx *fiber.Ctx) error {

	atID := ctx.Context().Value("token_at_id")
	if atID == nil {
		return fiber.NewError(fiber.StatusInternalServerError, "context value missing")
	}

	err := uc.userInteractor.SignOut(ctx.Context(), atID.(string))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "token invalidated, a new token is required to access the protected API\"",
	})
}

func (uc *userController) SignOutAll(ctx *fiber.Ctx) error {

	usrID := ctx.Context().Value("token_user_id")
	if usrID == nil {
		return fiber.NewError(fiber.StatusInternalServerError, "context value missing")
	}

	err := uc.userInteractor.SignOutAll(ctx.Context(), usrID.(string))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "token invalidated, a new token is required to access the protected API\"",
	})
}
