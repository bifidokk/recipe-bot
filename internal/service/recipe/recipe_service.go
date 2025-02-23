package recipe

import (
	"github.com/bifidokk/recipe-bot/internal/entity"
	"github.com/bifidokk/recipe-bot/internal/repository/recipe"
)

type Service interface {
	CreateRecipe(recipe *CreateRecipeData) (*entity.Recipe, error)
}

type recipeService struct {
	recipeRepository *recipe.Repository
}

func NewRecipeService(recipeRepository *recipe.Repository) Service {
	return &recipeService{
		recipeRepository: recipeRepository,
	}
}

func (r recipeService) CreateRecipe(_ *CreateRecipeData) (*entity.Recipe, error) {
	//TODO implement me
	panic("implement me")
}
