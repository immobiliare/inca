package middleware

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.rete.farm/sistemi/inca/server/helper"
)

func Bearer(acl map[string][]string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if len(acl) == 0 {
			return c.Next()
		}

		token := helper.GetToken(c, nil)
		if !helper.IsValidToken(token, acl) {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return c.Next()
	}
}
