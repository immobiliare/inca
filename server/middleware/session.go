package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/immobiliare/inca/server/helper"
)

const LoginPath = "/web/login"

func Session(store *session.Store, acl map[string][]string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Path() == LoginPath || len(acl) == 0 {
			_ = c.Bind(fiber.Map{"unprotected": true})
			return c.Next()
		}

		token := helper.GetToken(c, store)
		if !helper.IsValidToken(token, acl) {
			return c.Redirect(fmt.Sprintf("%s?redirect=%s", LoginPath, c.Path()), 302)
		}

		_ = c.Bind(fiber.Map{"authenticated": true})
		return c.Next()
	}
}
