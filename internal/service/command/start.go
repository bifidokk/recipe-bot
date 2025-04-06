package command

import (
	"github.com/rs/zerolog/log"
	"gopkg.in/telebot.v4"
)

const (
	startMenuText = "👨‍🍳 Hey there, food lover! Send me any TikTok recipe video and I'll instantly convert it into a clear, easy-to-follow recipe text. No more pausing or squinting at the screen! Just share the URL and let's start cooking! 🍽️"
)

type StartCommand struct{}

func NewStartCommand() Command {
	return &StartCommand{}
}

func (c *StartCommand) Name() string {
	return "start"
}

func (c *StartCommand) Register(b *telebot.Bot) {
	b.Handle("/start", func(ctx telebot.Context) error {
		log.Info().Msgf("/start command")

		menu := &telebot.ReplyMarkup{}
		btnRecipes := menu.Data("My recipes", "user_recipes")
		menu.Inline(
			menu.Row(btnRecipes),
		)

		return ctx.Send(startMenuText, menu)
	})
}
