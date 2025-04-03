package command

import (
	"github.com/bifidokk/recipe-bot/internal/entity"
	"github.com/bifidokk/recipe-bot/internal/service/recipe"
	"github.com/rs/zerolog/log"
	"gopkg.in/telebot.v4"
)

type UserRecipesCommand struct {
	recipeService recipe.Service
}

func NewUserRecipesCommand(recipeService recipe.Service) Command {
	return &UserRecipesCommand{
		recipeService: recipeService,
	}
}

func (c *UserRecipesCommand) Name() string {
	return "user_recipes"
}

func (c *UserRecipesCommand) Register(b *telebot.Bot) {
	b.Handle("/recipes", func(ctx telebot.Context) error {
		log.Info().Msgf("/recipes command")

		u := ctx.Get("user").(*entity.User)
		recipes, err := c.recipeService.GetRecipesByUserID(u.ID)

		if err != nil {
			log.Error().Err(err).Msg("Failed to get recipes")

			_, _ = b.Send(ctx.Sender(), "Sorry but I could not get your recipes")
			return err
		}

		if len(recipes) == 0 {
			_, _ = b.Send(ctx.Sender(), "You have no recipes yet")
			return nil
		}

		rcp := recipes[0]

		_, err = b.Send(
			ctx.Sender(),
			rcp.GetRecipeMarkdownView(),
			&telebot.SendOptions{
				ParseMode: telebot.ModeMarkdownV2,
			},
		)

		if err != nil {
			log.Error().Err(err).Msg("Failed to send message")
			return err
		}

		return nil
	})
}
