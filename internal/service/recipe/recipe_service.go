package recipe

import (
	"context"

	"github.com/bifidokk/recipe-bot/internal/entity"
	"github.com/bifidokk/recipe-bot/internal/repository/recipe"
)

type Service interface {
	CreateRecipe(recipe *CreateRecipeData, userID int) (*entity.Recipe, error)
}

type recipeService struct {
	recipeRepository *recipe.Repository
}

func NewRecipeService(recipeRepository *recipe.Repository) Service {
	return &recipeService{
		recipeRepository: recipeRepository,
	}
}

func (r recipeService) CreateRecipe(recipeData *CreateRecipeData, userID int) (*entity.Recipe, error) {
	rcp := &entity.Recipe{
		UserID:             userID,
		Title:              recipeData.Title,
		Body:               recipeData.Body,
		RecipeMarkdownText: recipeData.RecipeMarkdownText,
		Source:             recipeData.Source,
		SourceID:           recipeData.SourceID,
		SourceIDType:       recipeData.SourceIDType,
		AudioURL:           recipeData.AudioURL,
		ShareURL:           recipeData.ShareURL,
	}

	ctx := context.Background()
	recipeID, err := r.recipeRepository.CreateRecipe(ctx, rcp)

	if err != nil {
		return nil, err
	}

	return r.getRecipeByID(recipeID)
}

func (r recipeService) getRecipeByID(ID int) (*entity.Recipe, error) {
	ctx := context.Background()

	rcp, err := r.recipeRepository.FindByID(ctx, ID)

	if err != nil {
		return nil, err
	}

	return rcp, nil
}
