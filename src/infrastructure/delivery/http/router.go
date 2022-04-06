package http

import (
	"auth-project/src/interface/controller"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

const (
	APIv1 = "/api/v1"
)

type AuthStore struct {
	Rdb *redis.Client
}

func NewRouter(app *fiber.App, c controller.APIController) *fiber.App {

	authApi := app.Group(APIv1 + "/auth")

	authApi.Post("/authenticate", c.Auth.Authenticate)
	authApi.Post("/refresh", c.Auth.RefreshToken)

	qrCodeAuth := authApi.Group("/qr-code")

	qrCodeAuth.Post("/:qrCodeToken", authMiddleware(c), c.QrCodeAuth.CreateAuthTokenByAuthQrCode)
	qrCodeAuth.Get("/websocket", webSocketMiddleware(), websocket.New(c.QrCodeAuth.QrCodeAuthWebsocket))

	userApi := app.Group(APIv1 + "/users")

	userApi.Post("/sign-up", c.User.SignUp)
	userApi.Post("/sign-up/send-code", c.User.SignUpSendOTP)

	userApi.Post("/change-password", authMiddleware(c), c.User.ChangeMyPassword)

	userApi.Post("/reset-password/send-code", c.User.SendCodeForResetUserPassword)
	userApi.Post("/reset-password/verify-code", c.User.VerifyResetUserPasswordCode)

	userApi.Get("/my-profile", authMiddleware(c), c.User.GetMyProfile)

	userApi.Put("/myself/info", authMiddleware(c), c.User.UpdateMyselfInfo)
	userApi.Put("/myself/email", authMiddleware(c), c.User.UpdateMyselfEmail)
	userApi.Put("/myself/phone", authMiddleware(c), c.User.UpdateMyselfPhone)

	userApi.Post("/sign-out", authMiddleware(c), c.User.SignOut)
	userApi.Post("/sign-out/all", authMiddleware(c), c.User.SignOutAll)

	twoFactorAuthApi := app.Group(APIv1 + "/2fa")

	twoFactorAuthApi.Get("/google/qr-code", authMiddleware(c), c.TwoFactorAuth.GenerateGoogleTwoFactorAuthQrCode)

	twoFactorAuthApi.Post("/re-send", twoFactorAuthMiddleware(c), c.TwoFactorAuth.ReSendTwoFactorAuthCode)
	twoFactorAuthApi.Post("/verify", twoFactorAuthMiddleware(c), c.TwoFactorAuth.VerifyTwoFactorAuthCode)

	twoFactorAuthApi.Put("/set-up", authMiddleware(c), c.TwoFactorAuth.SetUpTwoFactorAuth)
	twoFactorAuthApi.Delete("/delete", authMiddleware(c), c.TwoFactorAuth.DeleteTwoFactorAuth)

	otpApi := app.Group(APIv1 + "/code")

	otpApi.Post("/send", authMiddleware(c), c.Token.Send2faCode)
	otpApi.Post("/target-send", authMiddleware(c), c.Token.SendTarget2faCode)

	return app
}
