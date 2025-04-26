package command

import (
	"github.com/bifidokk/recipe-bot/internal/entity"
	"github.com/bifidokk/recipe-bot/internal/service/recipe"
	"github.com/rs/zerolog/log"
	"gopkg.in/telebot.v4"
	"strconv"
)

type UserRecipeDetailsCommand struct {
	recipeService recipe.Service
}

func NewUserRecipeDetailsCommand(recipeService recipe.Service) Command {
	return &UserRecipeDetailsCommand{
		recipeService: recipeService,
	}
}

func (c *UserRecipeDetailsCommand) Name() string {
	return "user_recipe_details"
}

func (c *UserRecipeDetailsCommand) Register(b *telebot.Bot) {
	menu := &telebot.ReplyMarkup{}
	btnRecipeDetails := menu.Data("Recipe details", "user_recipe_details")
	b.Handle(&btnRecipeDetails, func(ctx telebot.Context) error {
		log.Info().Msgf("/user_recipe_detail inline button")

		return c.getUserRecipeDetails(ctx)
	})
}

func (c *UserRecipeDetailsCommand) getUserRecipeDetails(ctx telebot.Context) error {
	u := ctx.Get("user").(*entity.User)
	recipeID, err := strconv.Atoi(ctx.Callback().Data)

	if err != nil {
		log.Error().Err(err).Msg("Failed to convert recipe ID")
		return err
	}

	log.Info().Msgf("Recipe ID: %v", recipeID)

	rcp, err := c.recipeService.GetRecipeDetailsByIDForUser(recipeID, u.ID)
	if err != nil {
		log.Error().Msgf("Failed to get recipe details for user %v %v", u.ID, err)
		ctx.Send("Sorry but I could not get recipe details")
		return nil
	}

	menu := &telebot.ReplyMarkup{}
	btnRecipes := menu.Data("My recipes", "user_recipes")
	menu.Inline(
		menu.Row(btnRecipes),
	)

	ctx.Send(rcp.GetRecipeMarkdownView(), menu, &telebot.SendOptions{
		ParseMode:   telebot.ModeMarkdownV2,
		ReplyMarkup: menu,
	})

	return nil
}
