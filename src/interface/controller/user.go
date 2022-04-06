package controller

import (
	"auth-project/src/domain/model"
	"auth-project/src/usecase/interactor"
	"auth-project/tools"
	"github.com/gofiber/fiber/v2"
)

type userController struct {
	userInteractor interactor.UserInteractor
}

type UserController interface {
	SignUp(ctx *fiber.Ctx) error
	SignUpSendOTP(ctx *fiber.Ctx) error

	VerifyResetUserPasswordCode(ctx *fiber.Ctx) error
	SendCodeForResetUserPassword(ctx *fiber.Ctx) error

	GetMyProfile(ctx *fiber.Ctx) error

	ChangeMyPassword(ctx *fiber.Ctx) error
	UpdateMyselfInfo(ctx *fiber.Ctx) error
	UpdateMyselfPhone(ctx *fiber.Ctx) error
	UpdateMyselfEmail(ctx *fiber.Ctx) error

	SignOut(ctx *fiber.Ctx) error
	SignOutAll(ctx *fiber.Ctx) error
}

func NewUserController(ui interactor.UserInteractor) UserController {
	return &userController{ui}
}

// SignUp accepts 2fa and data for sign up, verify 2fa and activate user
func (uc *userController) SignUp(ctx *fiber.Ctx) error {

	var signUpReq model.SignUpReq
	err := ctx.BodyParser(&signUpReq)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = uc.userInteractor.SignUp(ctx.Context(), &signUpReq)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "OK",
	})
}

// SignUpSendOTP accepts login, send 2fa and create deactivate user
func (uc *userController) SignUpSendOTP(ctx *fiber.Ctx) error {

	var signUpSendOTPReq model.SignUpSend2faCodeReq
	err := ctx.BodyParser(&signUpSendOTPReq)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = uc.userInteractor.SignUpSendOTP(ctx.Context(), &signUpSendOTPReq)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "OK",
	})
}

// VerifyResetUserPasswordCode accepts 2fa reset password code, verify him and sent reset password access token
func (uc *userController) VerifyResetUserPasswordCode(ctx *fiber.Ctx) error {

	var userResetPasswordReq model.VerifyResetUserPassword2faСodeReq
	err := ctx.BodyParser(&userResetPasswordReq)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	userResetPasswordReq.NewPassword, err = tools.VerifyPassword(userResetPasswordReq.NewPassword)
	if err != nil {
		return err
	}

	err = uc.userInteractor.VerifyResetUserPasswordCode(ctx.Context(), &userResetPasswordReq)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "OK",
	})
}

// SendCodeForResetUserPassword accepts the user's email or phone, check his details, if they exist, send the 2fa code
func (uc *userController) SendCodeForResetUserPassword(ctx *fiber.Ctx) error {

	var send2faCodeForResetUserPasswordReq model.Send2faCodeForResetUserPasswordReq
	err := ctx.BodyParser(&send2faCodeForResetUserPasswordReq)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	resp, err := uc.userInteractor.SendCodeForResetUserPassword(ctx.Context(), &send2faCodeForResetUserPasswordReq)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)

}

func (uc *userController) GetMyProfile(ctx *fiber.Ctx) error {

	usrID, ok := ctx.Context().Value("token_user_id").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "context value type invalid")
	}

	resp, err := uc.userInteractor.GetMyProfileByID(ctx.Context(), usrID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// ChangeMyPassword takes 2fa password reset code and old password, verifies them and changes user password
func (uc *userController) ChangeMyPassword(ctx *fiber.Ctx) error {

	var reqData *model.UserChangePasswordReq
	err := ctx.BodyParser(&reqData)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	usrID, ok := ctx.Context().Value("token_user_id").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "context value type invalid")
	}

	err = uc.userInteractor.ChangeMyPassword(ctx.Context(), reqData, usrID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "OK",
	})
}

// UpdateMyselfInfo  takes the user's information, and updates it from the user
func (uc *userController) UpdateMyselfInfo(ctx *fiber.Ctx) error {

	var reqData *model.UserUpdateInfoReq
	err := ctx.BodyParser(&reqData)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	usrID, ok := ctx.Context().Value("token_user_id").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "context value type invalid")
	}

	resp, err := uc.userInteractor.UpdateUserInfoByID(ctx.Context(), &model.UserUpdateInfoData{FullName: reqData.FullName},
		usrID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// UpdateMyselfPhone takes the phone and 2fa code, checks the code and updates the user's phone
func (uc *userController) UpdateMyselfPhone(ctx *fiber.Ctx) error {

	var reqData *model.UserPhoneUpdateReq
	err := ctx.BodyParser(&reqData)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	usrID, ok := ctx.Context().Value("token_user_id").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "context value type invalid")
	}

	resp, err := uc.userInteractor.UpdateMyselfPhone(ctx.Context(), reqData, usrID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// UpdateMyselfEmail takes the email and 2fa code, checks the code and updates the user's email
func (uc *userController) UpdateMyselfEmail(ctx *fiber.Ctx) error {

	var reqData *model.UserEmailUpdateReq
	err := ctx.BodyParser(&reqData)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	usrID, ok := ctx.Context().Value("token_user_id").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "context value type invalid")
	}

	resp, err := uc.userInteractor.UpdateMyselfEmail(ctx.Context(), reqData, usrID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

//SignOut invalidates the user session in redis
func (uc *userController) SignOut(ctx *fiber.Ctx) error {

	sessionID, ok := ctx.Context().Value("token_session_id").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "context value type invalid")
	}

	err := uc.userInteractor.SignOut(ctx.Context(), sessionID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "token invalidated, a new token is required to access the protected API\"",
	})
}

// SignOutAll invalidates all user sessions in redis
func (uc *userController) SignOutAll(ctx *fiber.Ctx) error {

	userId, ok := ctx.Context().Value("token_user_id").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "context value type invalid")
	}

	err := uc.userInteractor.SignOutAll(ctx.Context(), userId)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "token invalidated, a new token is required to access the protected API",
	})
}
