package command

import (
	"github.com/bifidokk/recipe-bot/internal/entity"
	"github.com/bifidokk/recipe-bot/internal/service/recipe"
	"github.com/rs/zerolog/log"
	"gopkg.in/telebot.v4"

	"strconv"
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
		return c.getUserRecipes(ctx, b)
	})

	menu := &telebot.ReplyMarkup{}
	btnRecipes := menu.Data("My recipes", "user_recipes")
	b.Handle(&btnRecipes, func(ctx telebot.Context) error {
		log.Info().Msgf("/user_recipes inline button")
		return c.getUserRecipes(ctx, b)
	})
}

func (c *UserRecipesCommand) getUserRecipes(ctx telebot.Context, b *telebot.Bot) error {
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

	menu := &telebot.ReplyMarkup{}
	var rows []telebot.Row

	for _, rcp := range recipes {
		btnRecipe := menu.Data(rcp.Title, "user_recipe_details", strconv.Itoa(rcp.ID))
		rows = append(rows, menu.Row(btnRecipe))
	}

	menu.Inline(rows...)

	_, err = b.Send(
		ctx.Sender(),
		"Here are your recipes:",
		menu,
	)

	if err != nil {
		log.Error().Err(err).Msg("Failed to send message")
		return err
	}

	return nil
}
