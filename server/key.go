package server

import (
	"bytes"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (inca *Inca) handlerKey(c *fiber.Ctx) error {
	name := c.Params("name")
	if !inca.authorizedTarget(name, c) {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	_, data, err := (*inca.Storage).Get(name)
	if err == nil {
		if strings.EqualFold(c.Get("Accept", "text/plain"), "application/json") {
			return c.JSON(struct {
				Key string `json:"key"`
			}{string(data)})
		}
		return c.SendStream(bytes.NewReader(data), len(data))
	}

	log.Info().Str("name", name).Err(err).Msg("cached key not found")
	return c.SendStatus(fiber.StatusNotFound)
}
