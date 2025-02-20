package middleware

import (
	"errors"
	"strconv"

	"github.com/bifidokk/recipe-bot/internal/entity"
	"github.com/bifidokk/recipe-bot/internal/service"
	"github.com/rs/zerolog/log"
	tb "gopkg.in/telebot.v4"
)

func TgAuth(userService service.UserService) tb.MiddlewareFunc {
	return func(next tb.HandlerFunc) tb.HandlerFunc {
		return func(c tb.Context) error {
			sender := c.Sender()

			if sender == nil {
				log.Error().Msg("telegram sender is nil")
				return errors.New("telegram sender is nil")
			}

			telegramID := c.Sender().ID
			user, err := userService.GetUser(telegramID)

			if err != nil {
				log.Error().Err(err).Msg("failed to get user")
				return err
			}

			if user == nil {
				user = &entity.User{
					TelegramID: strconv.FormatInt(c.Sender().ID, 10),
					Name:       c.Sender().Username,
				}

				err = userService.CreateUser(user)

				if err != nil {
					log.Error().Err(err).Msg("failed to create user")
					return err
				}
			}

			c.Set("user", user)

			return next(c)
		}
	}
}
