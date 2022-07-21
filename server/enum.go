package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gitlab.rete.farm/sistemi/inca/pki"
)

func (inca *Inca) handlerEnum(c *fiber.Ctx) error {
	filter := c.Params("filter", ".*")
	results, err := (*inca.Storage).Find(filter)
	if err != nil {
		log.Error().Err(err).Msg("unable to enumerate certificates")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	crts := []Crt{}
	for _, result := range results {
		crt, err := pki.ParseBytes(result)
		if err != nil {
			log.Error().Err(err).Msg("unable to parse certificate")
		}
		crts = append(crts, EncodeCrt(crt))
	}

	return c.JSON(struct {
		Results []Crt `json:"results"`
	}{crts})
}
