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

	authApi.Get("authenticate/qr-code", c.Auth.GenerateQrCode)
	authApi.Get("authenticate/qr-codes/:token", c.Auth.LoginQrCode)
	authApi.Get("authenticate/qr-code/websocket/:token", webSocketMiddleware(), websocket.New(c.Auth.QrPolling))

	authApi.Post("/refresh", c.Auth.RefreshToken)

	userApi := app.Group(APIv1 + "/users")

	userApi.Post("/", c.User.CreateUser)
	userApi.Post("/change-password", authMiddleware(c), c.User.ChangeUserPassword)
	userApi.Put("/me", authMiddleware(c), c.User.UpdateYourself)
	userApi.Put("/my/pin", authMiddleware(c), c.User.UpdatePinYourself)
	userApi.Delete("/my/pin", authMiddleware(c), c.User.DeletePinYourself)
	userApi.Post("/sign-out", authMiddleware(c), c.User.SignOut)
	userApi.Post("/sign-out/all", authMiddleware(c), c.User.SignOutAll)

	return app
}
