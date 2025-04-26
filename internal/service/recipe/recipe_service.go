package recipe

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bifidokk/recipe-bot/internal/entity"
	"github.com/bifidokk/recipe-bot/internal/repository/recipe"
)

type Service interface {
	CreateRecipe(recipe *CreateRecipeData, userID int) (*entity.Recipe, error)
	UpdateRecipe(recipe *entity.Recipe) error
	GetRecipesByUserID(userID int) ([]*entity.Recipe, error)
	GetRecipeDetailsByIDForUser(recipeID int, userID int) (*entity.Recipe, error)
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
	}

	if recipeData.ShareURL != "" {
		rcp.ShareURL = sql.NullString{String: recipeData.ShareURL, Valid: true}
	}

	ctx := context.Background()
	recipeID, err := r.recipeRepository.CreateRecipe(ctx, rcp)

	if err != nil {
		return nil, err
	}

	return r.getRecipeByID(recipeID)
}

func (r recipeService) UpdateRecipe(recipe *entity.Recipe) error {
	ctx := context.Background()
	err := r.recipeRepository.UpdateRecipe(ctx, recipe)

	if err != nil {
		return err
	}

	return nil
}

func (r recipeService) GetRecipesByUserID(userID int) ([]*entity.Recipe, error) {
	ctx := context.Background()

	rcps, err := r.recipeRepository.FindByUserID(ctx, userID)

	if err != nil {
		return nil, err
	}

	return rcps, nil
}

func (r recipeService) GetRecipeDetailsByIDForUser(recipeID int, userID int) (*entity.Recipe, error) {
	rcp, err := r.getRecipeByID(recipeID)
	if err != nil {
		return nil, err
	}

	if rcp.UserID != userID {
		return nil, errors.New("user does not have access to this recipe")
	}

	return rcp, nil
}

func (r recipeService) getRecipeByID(ID int) (*entity.Recipe, error) {
	ctx := context.Background()

	rcp, err := r.recipeRepository.FindByID(ctx, ID)

	if err != nil {
		return nil, err
	}

	return rcp, nil
}
