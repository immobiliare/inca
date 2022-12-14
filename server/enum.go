package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/immobiliare/inca/pki"
	"github.com/immobiliare/inca/provider"
	"github.com/rs/zerolog/log"
)

func (inca *Inca) handlerEnum(c *fiber.Ctx) error {
	filter := c.Params("filter", ".*")
	results, err := (*inca.Storage).Find(filter)
	if err != nil {
		log.Error().Err(err).Msg("unable to enumerate certificates")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	crts := []JSONCrt{}
	for _, result := range results {
		crt, err := pki.ParseBytes(result)
		if err != nil {
			log.Error().Err(err).Msg("unable to parse certificate")
		}
		crts = append(crts, EncodeCrt(crt, provider.GetByTargetName(crt.Subject.CommonName, nil, inca.Providers)))
	}

	return c.JSON(struct {
		Results []JSONCrt `json:"results"`
	}{crts})
}
