package middleware

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	tb "gopkg.in/telebot.v4"
)

func Logger() tb.MiddlewareFunc {
	return func(next tb.HandlerFunc) tb.HandlerFunc {
		return func(c tb.Context) error {
			data, _ := json.MarshalIndent(c.Update(), "", "  ")
			log.Info().Msg(string(data))
			return next(c)
		}
	}
}
