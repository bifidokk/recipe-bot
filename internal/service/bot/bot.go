package bot

import (
	"github.com/bifidokk/recipe-bot/internal/service/command"

	"github.com/bifidokk/recipe-bot/internal/service"
	"github.com/rs/zerolog/log"
	"gopkg.in/telebot.v4"
)

type botService struct {
	bot      *telebot.Bot
	commands []command.Command
}

func NewBotService(
	bot *telebot.Bot,
	commands []command.Command,
) service.BotService {
	return &botService{
		bot,
		commands,
	}
}

func (bs *botService) Start() error {
	log.Info().Msg("Starting bot")

	err := bs.bot.SetCommands([]telebot.Command{
		{Text: "menu", Description: "Show the menu"},
		{Text: "recipes", Description: "Show your recipes"},
	})

	if err != nil {
		return err
	}

	for _, cmd := range bs.commands {
		cmd.Register(bs.bot)
	}

	log.Info().Msg("Bot started!")
	go func() {
		bs.bot.Start()
	}()

	return nil
}
