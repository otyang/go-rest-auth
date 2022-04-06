package http

import (
	"auth-project/src/interface/controller"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// allows you to perform functions where authorization is required
func authMiddleware(c controller.APIController) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		err := c.Auth.ValidateAccessToken(ctx)
		if err != nil {
			return err
		}
		return ctx.Next()
	}
}

// allows for two-factor authentication function
func twoFactorAuthMiddleware(c controller.APIController) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		err := c.Auth.ValidateTwoFactorAuthToken(ctx)
		if err != nil {
			return err
		}
		return ctx.Next()
	}
}

func webSocketMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("User-Agent", string(c.Request().Header.UserAgent()))
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}
