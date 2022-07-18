// SPDX-FileCopyrightText: 2020 Eder Sosa <eder.sohe@gmail.com>
// SPDX-License-Identifier: MIT
//
// Code below is a modified version of https://github.com/edersohe/zflogger:
// only addition is support for Fiber v2.
package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type fields struct {
	ID         string
	RemoteIP   string
	Host       string
	Method     string
	Path       string
	Protocol   string
	StatusCode int
	Latency    time.Duration
	Error      error
	Stack      []byte
}

func (req *fields) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("id", req.ID).
		Str("ip", req.RemoteIP).
		Str("host", req.Host).
		Str("method", req.Method).
		Str("path", req.Path).
		Str("protocol", req.Protocol).
		Int("status", req.StatusCode).
		Str("latency", fmt.Sprint(req.Latency)).
		Str("tag", "request")

	if req.Error != nil {
		e.Err(req.Error)
	}

	if req.Stack != nil {
		e.Bytes("stack", req.Stack)
	}
}

func Logger(log zerolog.Logger, filter func(*fiber.Ctx) bool) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		if filter != nil && filter(c) {
			return c.Next()
		}

		start := time.Now()

		rid := c.Get(fiber.HeaderXRequestID)
		if rid == "" {
			rid = uuid.New().String()
			c.Set(fiber.HeaderXRequestID, rid)
		}

		fields := &fields{
			ID:       rid,
			RemoteIP: c.IP(),
			Method:   c.Method(),
			Host:     c.Hostname(),
			Path:     c.Path(),
			Protocol: c.Protocol(),
		}

		defer func() {
			rvr := recover()
			if rvr != nil {
				err, ok := rvr.(error)
				if !ok {
					err = fmt.Errorf("%v", rvr)
				}

				fields.Error = err
				fields.Stack = debug.Stack()

				c.Status(http.StatusInternalServerError)
				if err := c.JSON(map[string]interface{}{
					"status": http.StatusText(http.StatusInternalServerError),
				}); err != nil {
					// damn linter wants me to error-check
					_ = err.Error()
				}

			}

			fields.StatusCode = c.Response().StatusCode()
			fields.Latency = time.Since(start)

			switch {
			case rvr != nil:
				log.Error().EmbedObject(fields).Msg("panic recover")
			case fields.StatusCode >= 500:
				log.Error().EmbedObject(fields).Msg("server error")
			case fields.StatusCode >= 400:
				log.Error().EmbedObject(fields).Msg("client error")
			case fields.StatusCode >= 300:
				log.Warn().EmbedObject(fields).Msg("redirect")
			case fields.StatusCode >= 200:
				log.Info().EmbedObject(fields).Msg("success")
			case fields.StatusCode >= 100:
				log.Info().EmbedObject(fields).Msg("informative")
			default:
				log.Warn().EmbedObject(fields).Msg("unknown status")
			}
		}()

		return c.Next()
	}
}
