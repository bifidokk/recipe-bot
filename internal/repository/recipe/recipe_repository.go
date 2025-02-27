package recipe

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/bifidokk/recipe-bot/internal/client"
	"github.com/bifidokk/recipe-bot/internal/entity"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/rs/zerolog/log"
)

type Repository struct {
	db         *client.DBClient
	sqlBuilder sq.StatementBuilderType
}

const (
	tableName = "recipes"

	idColumn           = "id"
	titleColumn        = "title"
	bodyColumn         = "body"
	markdownColumn     = "markdown"
	sourceColumn       = "source"
	sourceIDColumn     = "source_id"
	sourceIDTypeColumn = "source_id_type"
	audioURLColumn     = "audio_url"
	shareURLColumn     = "share_url"
	userIDColumn       = "user_id"
	createdAtColumn    = "created_at"
	updatedAtColumn    = "updated_at"
)

func NewRecipeRepository(db *client.DBClient) *Repository {
	return &Repository{
		db:         db,
		sqlBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *Repository) CreateRecipe(ctx context.Context, recipe *entity.Recipe) (int, error) {
	query, args, err := r.sqlBuilder.Insert(tableName).
		Columns(
			titleColumn,
			bodyColumn,
			markdownColumn,
			sourceColumn,
			sourceIDColumn,
			sourceIDTypeColumn,
			audioURLColumn,
			shareURLColumn,
			userIDColumn,
		).
		Values(
			recipe.Title,
			recipe.Body,
			recipe.RecipeMarkdownText,
			recipe.Source,
			recipe.SourceID,
			recipe.SourceIDType,
			recipe.AudioURL,
			recipe.ShareURL,
			recipe.UserID,
		).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		log.Error().Err(err).Msg("failed to build query")
		return 0, err
	}

	var recipeID int
	err = pgxscan.Get(ctx, r.db.Pool, &recipeID, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to create recipe: %w", err)
	}

	return recipeID, nil
}

func (r *Repository) FindByID(ctx context.Context, id int) (*entity.Recipe, error) {
	var recipe entity.Recipe
	query, args, err := r.sqlBuilder.Select(
		idColumn,
		titleColumn,
		bodyColumn,
		markdownColumn,
		sourceColumn,
		sourceIDColumn,
		sourceIDTypeColumn,
		audioURLColumn,
		shareURLColumn,
		userIDColumn,
		createdAtColumn,
		updatedAtColumn,
	).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		ToSql()

	if err != nil {
		log.Error().Err(err).Msg("failed to build query")
		return nil, err
	}

	err = pgxscan.Get(ctx, r.db.Pool, &recipe, query, args...)
	if err != nil {
		log.Info().Msgf("could not find recipe id %d", id)
		return nil, fmt.Errorf("failed to fetch recipe by ID: %w", err)
	}

	return &recipe, nil
}
