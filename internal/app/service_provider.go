package app

import (
	"github.com/bifidokk/recipe-bot/internal/config"
	"github.com/bifidokk/recipe-bot/internal/service"
	"github.com/bifidokk/recipe-bot/internal/service/api/openai"
	"github.com/bifidokk/recipe-bot/internal/service/api/tikhub"
	botService "github.com/bifidokk/recipe-bot/internal/service/bot"
	"github.com/bifidokk/recipe-bot/internal/service/video"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"net/http"
	"time"
)

type serviceProvider struct {
	tgBotConfig     config.TgBotConfig
	tikTokAPIConfig config.TikTokAPIConfig
	openAIAPIConfig config.OpenAIAPIConfig

	botService   service.BotService
	openAIClient service.OpenAIClient
	tikhubClient service.TikHubClient
	videoService service.VideoService
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

func (sp *serviceProvider) TikTokAPIConfig() config.TikTokAPIConfig {
	if sp.tikTokAPIConfig == nil {
		tikTokAPIConfig, err := config.NewTikTokAPIConfig()

		if err != nil {
			log.Fatalf("failed to get tik tok API config: %v", err)
		}

		sp.tikTokAPIConfig = tikTokAPIConfig
	}

	return sp.tikTokAPIConfig
}

func (sp *serviceProvider) OpenAIAPIConfig() config.OpenAIAPIConfig {
	if sp.openAIAPIConfig == nil {
		openAIAPIConfig, err := config.NewOpenAIAPIConfig()

		if err != nil {
			log.Fatalf("failed to get open ai API config: %v", err)
		}

		sp.openAIAPIConfig = openAIAPIConfig
	}

	return sp.openAIAPIConfig
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

		sp.botService = botService.NewBotService(
			bot,
			sp.TikTokAPIConfig().Token(),
			sp.OpenAIClient(),
			sp.VideoService(),
		)
	}

	return sp.botService
}

func (sp *serviceProvider) OpenAIClient() service.OpenAIClient {
	if sp.openAIClient == nil {
		sp.openAIClient = openai.NewOpenAIClient(sp.OpenAIAPIConfig().Token())
	}

	return sp.openAIClient
}

func (sp *serviceProvider) TikhubClient() service.TikHubClient {
	if sp.tikhubClient == nil {
		sp.tikhubClient = tikhub.NewTikHubClient(sp.TikTokAPIConfig().Token())
	}

	return sp.tikhubClient
}

func (sp *serviceProvider) VideoService() service.VideoService {
	if sp.videoService == nil {
		sp.videoService = video.NewVideoService(sp.TikhubClient())
	}

	return sp.videoService
}
