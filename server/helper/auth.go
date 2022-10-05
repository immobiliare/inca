package helper

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func GetToken(c *fiber.Ctx, store *session.Store) string {
	var (
		header = c.Get("authorization")
		fields = strings.SplitN(header, " ", 2)
	)
	if len(fields) == 2 && strings.EqualFold(fields[0], "bearer") {
		return fields[1]
	}

	if store == nil {
		return ""
	}

	session, err := store.Get(c)
	if err != nil {
		return ""
	}

	name := session.Get("name")
	if name == nil {
		return ""
	}

	return name.(string)
}

func IsValidToken(token string, acl map[string][]string) bool {
	for key := range acl {
		if strings.EqualFold(key, token) {
			return true
		}
	}
	return false
}
