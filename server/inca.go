package server

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/redirect/v2"
	"github.com/gofiber/template/django"
	"github.com/immobiliare/inca/provider"
	"github.com/immobiliare/inca/server/config"
	"github.com/immobiliare/inca/server/helper"
	"github.com/immobiliare/inca/server/middleware"
	"github.com/immobiliare/inca/storage"
	"github.com/immobiliare/inca/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Inca struct {
	*fiber.App
	Storage      *storage.Storage
	Providers    []*provider.Provider
	sessionStore *session.Store
	acl          map[string][]string
}

func Spinup(path string) (*Inca, error) {
	cfg, err := config.Parse(path)
	if err != nil {
		return nil, err
	}

	if len(cfg.Sentry) > 0 {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              cfg.Sentry,
			TracesSampleRate: 1.0,
		}); err != nil {
			return nil, fmt.Errorf("sentry: %s", err)
		}
		defer sentry.Flush(2 * time.Second)
		log.Info().Msg("sentry correctly initialized")
	}

	templateEngine := django.New(cfg.TemplatesPath, ".html.j2")
	templateEngine.Reload(strings.EqualFold(cfg.Environment, "development"))
	templateEngine.Debug(strings.EqualFold(cfg.Environment, "development"))

	inca := &Inca{
		fiber.New(
			fiber.Config{
				DisableStartupMessage: true,
				Views:                 templateEngine,
				// Views:                 html.NewFileSystem(http.Dir("./server/views"), ".html.j2"),
			},
		),
		cfg.Storage,
		cfg.Providers,
		session.New(),
		cfg.ACL,
	}
	inca.Use(compress.New())
	inca.Use(middleware.Logger(zerolog.New(os.Stdout), func(c *fiber.Ctx) bool {
		return strings.HasPrefix(c.Path(), "/health") ||
			strings.HasPrefix(c.Path(), "/static/") ||
			strings.HasPrefix(c.Path(), "/favicon.ico")
	}))
	inca.Use(redirect.New(redirect.Config{
		Rules: map[string]string{
			"^/$":            "/web",
			"^/favicon.ico$": "/static/favicon.ico",
		},
		StatusCode: 301,
	}))

	static := fiber.Static{
		Compress:      true,
		CacheDuration: 24 * time.Hour,
	}
	if strings.EqualFold(cfg.Environment, "development") {
		static.Compress = false
		static.CacheDuration = 5 * time.Second
	}
	inca.Static("/static", "./server/static", static)
	incaWeb := inca.Group("/web")
	incaWeb.Use(middleware.Session(inca.sessionStore, inca.acl))
	incaWeb.Get("/", inca.handlerWebIndex)
	incaWeb.Get("/login", inca.handlerWebAuthLoginView)
	incaWeb.Post("/login", inca.handlerWebAuthLogin)
	incaWeb.Get("/logout", inca.handlerWebAuthLogout)
	incaWeb.Get("/config", inca.handlerWebConfig)
	incaWeb.Get("/issue", inca.handlerWebIssueView)
	incaWeb.Post("/issue", inca.handlerWebIssue)
	incaWeb.Get("/import", inca.handlerWebImportView)
	incaWeb.Post("/import", inca.handlerWebImport)
	incaWeb.Get("/:name", inca.handlerWebView)
	incaWeb.Post("/:name/delete", inca.handlerWebDelete)
	inca.Get("/enum", inca.handlerEnum)
	inca.Get("/health", inca.handlerHealth)
	inca.Get("/ca/:filter", inca.handlerCA)
	inca.Get("/:name", inca.handlerCRT)
	inca.Get("/:name/key", inca.handlerKey)
	inca.Get("/:name/show", inca.handlerShow)
	inca.Delete("/:name", inca.handlerRevoke)
	return inca, nil
}

func (inca *Inca) authorizedTarget(name string, c *fiber.Ctx) bool {
	if len(inca.acl) == 0 {
		return true
	}
	if targets, ok := inca.acl[helper.GetToken(c, inca.sessionStore)]; ok {
		return util.RegexesAnyMatch(name, targets...)
	}
	return false
}
