package server

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gitlab.rete.farm/sistemi/inca/util"
)

func (inca *Inca) handlerKey(c *fiber.Ctx) error {
	_, data, err := (*inca.Cfg.Storage).Get(c.Params("name"))
	if err == nil {
		if strings.EqualFold(c.Get("Accept", "text/plain"), "application/json") {
			return c.SendString(fmt.Sprintf(`{"key":"%s"}`, util.BytesToJSON(data)))
		}
		return c.SendStream(bytes.NewReader(data), len(data))
	}

	log.Info().Str("name", c.Params("name")).Err(err).Msg("cached key not found")
	return c.SendStatus(fiber.StatusNotFound)
}
