package http

import (
	"auth-project/src/interface/controller"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)


func authMiddleware(c controller.APIController) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		err := c.Auth.ValidateAccessTokenID(ctx)
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
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}


