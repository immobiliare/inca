package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gitlab.rete.farm/sistemi/inca/server/helper"
)

const (
	invalidSessionName = "invalid"
	LoginPath          = "/web/login"
)

func Session(store *session.Store, acl map[string][]string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Path() == LoginPath {
			return c.Next()
		}

		if !IsAuthenticated(c, store, acl) {
			return c.Redirect(fmt.Sprintf("%s?redirect=%s", LoginPath, c.Path()), 302)
		}

		_ = c.Bind(fiber.Map{"authenticated": true})
		return c.Next()
	}
}

func IsAuthenticated(c *fiber.Ctx, store *session.Store, acl map[string][]string) bool {
	session, err := store.Get(c)
	if err != nil {
		return false
	}

	if name := session.Get("name"); name == nil ||
		len(name.(string)) == 0 ||
		strings.EqualFold(name.(string), invalidSessionName) ||
		!helper.IsValidToken(name.(string), acl) {
		return false
	}
	return true
}
