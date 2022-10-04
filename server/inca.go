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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.rete.farm/sistemi/inca/provider"
	"gitlab.rete.farm/sistemi/inca/server/config"
	"gitlab.rete.farm/sistemi/inca/server/middleware"
	"gitlab.rete.farm/sistemi/inca/storage"
)

type Inca struct {
	*fiber.App
	Storage      *storage.Storage
	Providers    []*provider.Provider
	sessionStore *session.Store
	ACL          map[string][]string
}

type ACL map[string][]string

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
	incaWeb.Use(middleware.Session(inca.sessionStore, inca.ACL))
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
	incaAPI := inca.Group("/")
	incaAPI.Use(middleware.Bearer(inca.ACL))
	incaAPI.Get("/enum", inca.handlerEnum)
	incaAPI.Get("/health", inca.handlerHealth)
	incaAPI.Get("/ca/:filter", inca.handlerCA)
	incaAPI.Get("/:name", inca.handlerCRT)
	incaAPI.Get("/:name/key", inca.handlerKey)
	incaAPI.Get("/:name/show", inca.handlerShow)
	incaAPI.Delete("/:name", inca.handlerRevoke)
	return inca, nil
}
