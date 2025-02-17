package app

import (
	"net/http"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/bifidokk/recipe-bot/internal/config"
	"github.com/bifidokk/recipe-bot/internal/service"
	"github.com/bifidokk/recipe-bot/internal/service/api/openai"
	"github.com/bifidokk/recipe-bot/internal/service/api/tikhub"
	botService "github.com/bifidokk/recipe-bot/internal/service/bot"
	"github.com/bifidokk/recipe-bot/internal/service/video"
	"github.com/rs/zerolog/log"
	tb "gopkg.in/tucnak/telebot.v2"
)

type serviceProvider struct {
	tgBotConfig     config.TgBotConfig
	tikTokAPIConfig config.TikTokAPIConfig
	openAIAPIConfig config.OpenAIAPIConfig
	pgConfig        config.PgConfig

	db *gorm.DB

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
			log.Fatal().Err(err).Msg("Failed to get tg bot config")
		}

		sp.tgBotConfig = tgBotConfig
	}

	return sp.tgBotConfig
}

func (sp *serviceProvider) TikTokAPIConfig() config.TikTokAPIConfig {
	if sp.tikTokAPIConfig == nil {
		tikTokAPIConfig, err := config.NewTikTokAPIConfig()

		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get tik tok API config")
		}

		sp.tikTokAPIConfig = tikTokAPIConfig
	}

	return sp.tikTokAPIConfig
}

func (sp *serviceProvider) OpenAIAPIConfig() config.OpenAIAPIConfig {
	if sp.openAIAPIConfig == nil {
		openAIAPIConfig, err := config.NewOpenAIAPIConfig()

		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get open ai API config")
		}

		sp.openAIAPIConfig = openAIAPIConfig
	}

	return sp.openAIAPIConfig
}

func (sp *serviceProvider) PgConfig() config.PgConfig {
	if sp.pgConfig == nil {
		pgConfig, err := config.NewPgConfig()

		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get pg config")
		}

		sp.pgConfig = pgConfig
	}

	return sp.pgConfig
}

func (sp *serviceProvider) DB() *gorm.DB {
	if sp.db == nil {
		db, err := gorm.Open(postgres.Open(sp.PgConfig().Dsn()), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to database")
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get database instance")
		}

		log.Info().Msg("Connected to database")

		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(5 * time.Minute)

		sp.db = db
	}

	return sp.db
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
			log.Fatal().Err(err).Msg("Failed  to init tg bot")
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
