package app

import (
	"context"
	"github.com/bifidokk/recipe-bot/internal/config"
	tb "gopkg.in/tucnak/telebot.v2"
)

type App struct {
	serviceProvider *serviceProvider
	bot             *tb.Bot
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
	}

	for _, initFunction := range inits {
		if err := initFunction(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) initConfig(ctx context.Context) error {
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
