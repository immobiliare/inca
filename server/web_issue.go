package server

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gitlab.rete.farm/sistemi/inca/pki"
	"gitlab.rete.farm/sistemi/inca/provider"
	"gitlab.rete.farm/sistemi/inca/util"
)

func (inca *Inca) handlerWebIssueView(c *fiber.Ctx) error {
	return c.Render("issue", fiber.Map{
		"algorithms": util.StringSliceDistinct(
			append([]string{pki.DefaultCrtAlgo}, pki.EDDSA, pki.ECDSA, pki.RSA)),
	})
}

func (inca *Inca) handlerWebIssue(c *fiber.Ctx) error {
	var (
		options = util.ParseQueryString(c.Body())
		name    string
	)
	if names, ok := options["alt"]; !ok || len(names) == 0 {
		_ = c.Bind(fiber.Map{"error": "At least a certificate name gotta be given"})
		return inca.handlerWebIssueView(c)
	} else {
		name = strings.Split(names, ",")[0]
	}

	if !inca.authorizedTarget(name, c) {
		_ = c.Bind(fiber.Map{"error": "Unauthorized to issue the certificate"})
		return inca.handlerWebIssueView(c)
	}

	if !pki.IsValidCN(name) {
		_ = c.Bind(fiber.Map{"error": "Invalid certificate name"})
		return inca.handlerWebIssueView(c)
	}

	p := provider.GetByTargetName(name, options, inca.Providers)
	if p == nil {
		_ = c.Bind(fiber.Map{"error": "Unable to find a suitable provider"})
		return inca.handlerWebIssueView(c)
	}

	if _, _, err := (*inca.Storage).Get(name); err == nil {
		_ = c.Bind(fiber.Map{"error": "Certificate already existing"})
		return inca.handlerWebIssueView(c)
	}

	crt, key, err := (*p).Get(name, options)
	if err != nil {
		log.Error().Err(err).Msg("unable to issue certificate")
		_ = c.Bind(fiber.Map{"error": "Unable to issue the certificate"})
		return inca.handlerWebIssueView(c)
	}

	if err := (*inca.Storage).Put(name, crt, key); err != nil {
		log.Error().Err(err).Msg("unable to persist certificate")
		_ = c.Bind(fiber.Map{"error": "Unable to persist the certificate"})
		return inca.handlerWebIssueView(c)
	}

	_ = c.Bind(fiber.Map{"message": "Certificate successfully created"})
	return inca.handlerWebIndex(c)
}
