package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"gitlab.rete.farm/sistemi/inca/server/helper"
)

func Bearer(acl map[string][]string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorization := c.Get("authorization")
		if len(authorization) == 0 {
			return c.SendStatus(400)
		}

		var (
			fields = strings.SplitN(authorization, " ", 2)
			label  = fields[0]
			token  = fields[1]
		)
		if !(strings.EqualFold(label, "bearer") && helper.IsValidToken(token, acl)) {
			return c.SendStatus(400)
		}

		return c.Next()
	}
}
