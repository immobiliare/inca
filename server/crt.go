package server

import (
	"bytes"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/immobiliare/inca/pki"
	"github.com/immobiliare/inca/provider"
	"github.com/immobiliare/inca/util"
	"github.com/rs/zerolog/log"
)

func (inca *Inca) handlerCRT(c *fiber.Ctx) error {
	var (
		name         = c.Params("name")
		queryStrings = util.ParseQueryString(c.Request().URI().QueryString())
		alt          = append([]string{name}, queryStringAlt(queryStrings)...)
	)
	if !inca.authorizedTarget(name, c) {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	if !pki.IsValidCN(name) {
		log.Error().Str("name", name).Msg("invalid name")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	data, _, err := (*inca.Storage).Get(name)
	if err == nil {
		if crt, err := pki.ParseBytes(data); err != nil {
			log.Error().Str("name", name).Msg("unable to parse cached certificate")
		} else {
			crtDNSNames, crtIPAddresses := pki.AltNames(crt)
			reqDNSNames, reqIPAddresses := pki.ParseAltNames(alt)
			dnsNames, ipAddresses := util.StringSliceDistinct(append(crtDNSNames, reqDNSNames...)),
				util.StringSliceDistinct(append(crtIPAddresses, reqIPAddresses...))
			if !util.StringSlicesEqual(crtDNSNames, dnsNames) ||
				!util.StringSlicesEqual(crtIPAddresses, ipAddresses) {
				log.Info().Str("name", name).Msg("cached certificate needs flush")
				queryStrings["alt"] = strings.Join(append(dnsNames, ipAddresses...), ",")
			} else if crt.NotAfter.Before(time.Now()) {
				log.Info().Str("name", name).Msg("certificate is expired: going to renew it")
			} else {
				log.Info().Str("name", name).Msg("returning cached certificate")
				if strings.EqualFold(c.Get("Accept", "text/plain"), "application/json") {
					return c.JSON(struct {
						Crt string `json:"crt"`
					}{string(data)})
				}
				return c.SendStream(bytes.NewReader(data), len(data))
			}
		}
	}

	p := provider.GetByTargetName(name, queryStrings, inca.Providers)
	if p == nil {
		log.Error().Str("name", name).Msg("no provider found")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	crt, key, err := (*p).Get(name, queryStrings)
	if err != nil {
		log.Error().Err(err).Msg("unable to generate")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if err := (*inca.Storage).Put(name, crt, key); err != nil {
		log.Error().Err(err).Msg("unable to persist certificate")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if strings.EqualFold(c.Get("Accept", "text/plain"), "application/json") {
		return c.JSON(struct {
			Crt string `json:"crt"`
		}{string(crt)})
	}
	return c.SendStream(bytes.NewReader(crt), len(crt))
}

func queryStringAlt(queryStrings map[string]string) (altNames []string) {
	if param, ok := queryStrings["alt"]; ok {
		altNames = strings.Split(param, ",")
	}
	return
}
