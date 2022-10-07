package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/immobiliare/inca/server/helper"
	"github.com/immobiliare/inca/server/middleware"
	"github.com/immobiliare/inca/util"
	"github.com/rs/zerolog/log"
)

func (inca *Inca) handlerWebAuthLogin(c *fiber.Ctx) error {
	var (
		options  = util.ParseQueryString(c.Body())
		redirect = c.Query("redirect", "/web")
	)
	token, ok := options["token"]
	if !ok {
		_ = c.Bind(fiber.Map{"error": "No token is given"})
		return inca.handlerWebAuthLoginView(c)
	} else if !helper.IsValidToken(token, inca.acl) {
		_ = c.Bind(fiber.Map{"error": "Unauthorized token"})
		return inca.handlerWebAuthLoginView(c)

	}

	session, err := inca.sessionStore.Get(c)
	if err != nil {
		_ = c.Bind(fiber.Map{"error": "Unable to spawn a session"})
		log.Error().Err(err).Msg("unable to spawn a session")
		return inca.handlerWebAuthLoginView(c)
	}

	session.Set("name", token)
	if err := session.Save(); err != nil {
		_ = c.Bind(fiber.Map{"error": "Unable to persist the session"})
		log.Error().Err(err).Msg("unable to persist the session")
		return inca.handlerWebAuthLoginView(c)
	}

	return c.Redirect(redirect)
}

func (inca *Inca) handlerWebAuthLoginView(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{})
}

func (inca *Inca) handlerWebAuthLogout(c *fiber.Ctx) error {
	session, err := inca.sessionStore.Get(c)
	if err != nil {
		_ = c.Bind(fiber.Map{"error": "Unable to find the session"})
		log.Error().Err(err).Msg("unable to find the session")
		return inca.handlerWebIndex(c)
	}

	session.Delete("name")
	if err := session.Destroy(); err != nil {
		_ = c.Bind(fiber.Map{"error": "Unable to destroy the session"})
		log.Error().Err(err).Msg("unable to destroy the session")
		return inca.handlerWebIndex(c)
	}

	return c.Redirect(middleware.LoginPath)
}
