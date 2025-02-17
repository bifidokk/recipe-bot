package app

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/bifidokk/recipe-bot/internal/config"
)

type App struct {
	serviceProvider *serviceProvider
}

func NewApp(ctx context.Context) (*App, error) {
	application := &App{}
	err := application.initDependencies(ctx)

	if err != nil {
		return nil, err
	}

	return application, nil
}

func (app *App) Run() error {
	err := app.serviceProvider.BotService().Start()

	return err
}

func (app *App) initDependencies(ctx context.Context) error {
	inits := []func(context context.Context) error{
		app.initConfig,
		app.initServiceProvider,
		app.initLogger,
	}

	for _, initFunction := range inits {
		if err := initFunction(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

func (app *App) initServiceProvider(_ context.Context) error {
	app.serviceProvider = newServiceProvider()

	return nil
}

func (app *App) initLogger(_ context.Context) error {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return err
	}

	multi := io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stdout}, logFile)
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()

	go func(logFile *os.File, interval time.Duration) {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			err := logFile.Sync()
			if err != nil {
				log.Error().Err(err).Msg("Failed to sync log file")
			}
		}
	}(logFile, 60*time.Second)

	return nil
}
