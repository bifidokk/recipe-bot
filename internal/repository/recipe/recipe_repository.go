package recipe

import (
	"context"
	"fmt"
	"time"

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
)

type columns struct {
	id           string
	title        string
	body         string
	markdown     string
	source       string
	sourceID     string
	sourceIDType string
	audioURL     string
	shareURL     string
	userID       string
	coverFileID  string
	createdAt    string
	updatedAt    string
}

var cols = columns{
	id:           "id",
	title:        "title",
	body:         "body",
	markdown:     "markdown",
	source:       "source",
	sourceID:     "source_id",
	sourceIDType: "source_id_type",
	audioURL:     "audio_url",
	shareURL:     "share_url",
	userID:       "user_id",
	coverFileID:  "cover_file_id",
	createdAt:    "created_at",
	updatedAt:    "updated_at",
}

func (c columns) all() []string {
	return []string{
		c.id,
		c.title,
		c.body,
		c.markdown,
		c.source,
		c.sourceID,
		c.sourceIDType,
		c.audioURL,
		c.shareURL,
		c.userID,
		c.coverFileID,
		c.createdAt,
		c.updatedAt,
	}
}

func (c columns) forInsert() []string {
	return []string{
		c.title,
		c.body,
		c.markdown,
		c.source,
		c.sourceID,
		c.sourceIDType,
		c.audioURL,
		c.shareURL,
		c.coverFileID,
		c.userID,
	}
}

func NewRecipeRepository(db *client.DBClient) *Repository {
	return &Repository{
		db:         db,
		sqlBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *Repository) CreateRecipe(ctx context.Context, recipe *entity.Recipe) (int, error) {
	query, args, err := r.sqlBuilder.Insert(tableName).
		Columns(cols.forInsert()...).
		Values(
			recipe.Title,
			recipe.Body,
			recipe.RecipeMarkdownText,
			recipe.Source,
			recipe.SourceID,
			recipe.SourceIDType,
			recipe.AudioURL,
			recipe.ShareURL,
			recipe.CoverFileID,
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
	query, args, err := r.sqlBuilder.Select(cols.all()...).
		From(tableName).
		Where(sq.Eq{cols.id: id}).
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

func (r *Repository) UpdateRecipe(ctx context.Context, recipe *entity.Recipe) error {
	updates := map[string]interface{}{
		cols.title:        recipe.Title,
		cols.body:         recipe.Body,
		cols.markdown:     recipe.RecipeMarkdownText,
		cols.source:       recipe.Source,
		cols.sourceID:     recipe.SourceID,
		cols.sourceIDType: recipe.SourceIDType,
		cols.audioURL:     recipe.AudioURL,
		cols.shareURL:     recipe.ShareURL,
		cols.coverFileID:  recipe.CoverFileID,
		cols.userID:       recipe.UserID,
	}

	updateBuilder := r.sqlBuilder.Update(tableName).Where("id = ?", recipe.ID)

	for column, value := range updates {
		updateBuilder = updateBuilder.Set(column, value)
	}

	updateBuilder = updateBuilder.Set("updated_at", time.Now())

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		log.Error().Err(err).Msg("failed to build update query")
		return err
	}

	_, err = r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update recipe: %w", err)
	}

	return nil
}

func (r *Repository) FindByUserID(ctx context.Context, userID int) ([]*entity.Recipe, error) {
	var recipes []*entity.Recipe
	query, args, err := r.sqlBuilder.Select(cols.all()...).
		From(tableName).
		Where(sq.Eq{cols.userID: userID}).
		ToSql()

	if err != nil {
		log.Error().Err(err).Msg("failed to build query")
		return nil, err
	}

	err = pgxscan.Select(ctx, r.db.Pool, &recipes, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch recipes by user ID: %w", err)
	}

	return recipes, nil
}
