package app

import (
	"net/http"
	"time"

	recipeRepo "github.com/bifidokk/recipe-bot/internal/repository/recipe"
	userRepo "github.com/bifidokk/recipe-bot/internal/repository/user"
	"github.com/bifidokk/recipe-bot/internal/service/recipe"

	"github.com/bifidokk/recipe-bot/internal/service/user"

	"github.com/bifidokk/recipe-bot/internal/middleware"

	"github.com/bifidokk/recipe-bot/internal/client"

	"github.com/bifidokk/recipe-bot/internal/config"
	"github.com/bifidokk/recipe-bot/internal/service"
	"github.com/bifidokk/recipe-bot/internal/service/api/instaloader"
	"github.com/bifidokk/recipe-bot/internal/service/api/openai"
	"github.com/bifidokk/recipe-bot/internal/service/api/tikhub"
	botService "github.com/bifidokk/recipe-bot/internal/service/bot"
	"github.com/bifidokk/recipe-bot/internal/service/command"
	"github.com/bifidokk/recipe-bot/internal/service/video"
	"github.com/rs/zerolog/log"
	tb "gopkg.in/telebot.v4"
)

type serviceProvider struct {
	tgBotConfig     config.TgBotConfig
	tikTokAPIConfig config.TikTokAPIConfig
	openAIAPIConfig config.OpenAIAPIConfig
	pgConfig        config.PgConfig

	db *client.DBClient

	botService        service.BotService
	openAIClient      service.OpenAIClient
	tikhubClient      service.TikHubClient
	instaloaderClient service.InstaloaderClient
	videoService      service.VideoService
	userService       service.UserService
	recipeService     recipe.Service

	userRepository   *userRepo.Repository
	recipeRepository *recipeRepo.Repository

	botCommands []command.Command
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

func (sp *serviceProvider) DB() *client.DBClient {
	if sp.db == nil {
		db, err := client.NewDBClient(sp.PgConfig().Dsn())

		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to database")
		}

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

		bot.Use(middleware.Logger())
		bot.Use(middleware.TgAuth(sp.UserService()))

		sp.botService = botService.NewBotService(
			bot,
			sp.BotCommands(),
		)
	}

	return sp.botService
}

func (sp *serviceProvider) BotCommands() []command.Command {
	if sp.botCommands == nil {
		sp.botCommands = []command.Command{
			command.NewStartCommand(),
			command.NewUserRecipesCommand(sp.RecipeService()),
			command.NewTextCommand(
				sp.OpenAIClient(),
				sp.VideoService(),
				sp.RecipeService(),
				sp.UserService(),
			),
			command.NewUserRecipeDetailsCommand(sp.RecipeService()),
		}
	}

	return sp.botCommands
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

func (sp *serviceProvider) InstaloaderClient() service.InstaloaderClient {
	if sp.instaloaderClient == nil {
		sp.instaloaderClient = instaloader.NewInstaloaderClient()
	}

	return sp.instaloaderClient
}

func (sp *serviceProvider) VideoService() service.VideoService {
	if sp.videoService == nil {
		sp.videoService = video.NewVideoService(
			sp.TikhubClient(),
			sp.InstaloaderClient(),
		)
	}

	return sp.videoService
}

func (sp *serviceProvider) UserRepository() *userRepo.Repository {
	if sp.userRepository == nil {
		sp.userRepository = userRepo.NewUserRepository(sp.DB())
	}

	return sp.userRepository
}

func (sp *serviceProvider) UserService() service.UserService {
	if sp.userService == nil {
		sp.userService = user.NewUserService(sp.UserRepository())
	}

	return sp.userService
}

func (sp *serviceProvider) RecipeService() recipe.Service {
	if sp.recipeService == nil {
		sp.recipeService = recipe.NewRecipeService(sp.RecipeRepository())
	}

	return sp.recipeService
}

func (sp *serviceProvider) RecipeRepository() *recipeRepo.Repository {
	if sp.recipeRepository == nil {
		sp.recipeRepository = recipeRepo.NewRecipeRepository(sp.DB())
	}

	return sp.recipeRepository
}
