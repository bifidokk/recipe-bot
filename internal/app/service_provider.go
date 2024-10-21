package app

import (
	"github.com/bifidokk/receipt-bot/internal/config"
	"github.com/bifidokk/receipt-bot/internal/service"
	botService "github.com/bifidokk/receipt-bot/internal/service/bot"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"net/http"
	"time"
)

type serviceProvider struct {
	tgBotConfig config.TgBotConfig

	botService service.BotService
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (sp *serviceProvider) TgBotConfig() config.TgBotConfig {
	if sp.tgBotConfig == nil {
		tgBotConfig, err := config.NewTgBotConfig()

		if err != nil {
			log.Fatalf("failed to get tg bot config: %v", err)
		}

		sp.tgBotConfig = tgBotConfig
	}

	return sp.tgBotConfig
}

func (sp *serviceProvider) BotService() service.BotService {
	if sp.botService == nil {
		var httpClient *http.Client

		bot, err := tb.NewBot(tb.Settings{
			Token:  sp.TgBotConfig().Token(),
			Poller: &tb.LongPoller{Timeout: 10 * time.Second},
			Client: httpClient,
		})

		if err != nil {
			log.Fatalf("failed to init tg bot: %v", err)
		}

		sp.botService = botService.NewBotService(bot)
	}

	return sp.botService
}
