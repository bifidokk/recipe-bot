package user

import (
	"context"
	"fmt"
	"strconv"

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

type columns struct {
	id          string
	name        string
	tgID        string
	language    string
	recipeLimit string
	createdAt   string
	updatedAt   string
}

var cols = columns{
	id:          "id",
	name:        "name",
	tgID:        "telegram_id",
	language:    "language_code",
	recipeLimit: "recipe_limit",
	createdAt:   "created_at",
	updatedAt:   "updated_at",
}

func (c columns) all() []string {
	return []string{
		c.id,
		c.name,
		c.tgID,
		c.language,
		c.recipeLimit,
		c.createdAt,
		c.updatedAt,
	}
}

func (c columns) forInsert() []string {
	return []string{
		c.name,
		c.tgID,
		c.language,
		c.recipeLimit,
	}
}

const (
	tableName = "users"
)

func NewUserRepository(dbClient *client.DBClient) *Repository {
	return &Repository{
		db:         dbClient,
		sqlBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *Repository) FindByTelegramID(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User
	query, args, err := r.sqlBuilder.Select(
		cols.all()...,
	).
		From(tableName).
		Where(sq.Eq{cols.tgID: strconv.Itoa(int(id))}).
		ToSql()

	if err != nil {
		log.Error().Err(err).Msg("failed to build query")
		return nil, err
	}

	err = pgxscan.Get(ctx, r.db.Pool, &user, query, args...)
	if err != nil {
		log.Info().Msgf("could not find user with tg id %d", id)
		return nil, nil
	}

	return &user, nil
}

func (r *Repository) FindByID(ctx context.Context, id int) (*entity.User, error) {
	var user entity.User
	query, args, err := r.sqlBuilder.Select(
		cols.all()...,
	).
		From(tableName).
		Where(sq.Eq{cols.id: id}).
		ToSql()

	if err != nil {
		log.Error().Err(err).Msg("failed to build query")
		return nil, err
	}

	err = pgxscan.Get(ctx, r.db.Pool, &user, query, args...)
	if err != nil {
		log.Info().Msgf("could not find user id %d", id)
		return nil, fmt.Errorf("failed to fetch user by ID: %w", err)
	}

	return &user, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *entity.User) (int, error) {
	query, args, err := r.sqlBuilder.Insert(tableName).
		Columns(cols.forInsert()...).
		Values(user.Name, user.TelegramID, user.Language, user.RecipeLimit).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		log.Error().Err(err).Msg("failed to build query")
		return 0, err
	}

	var userID int
	err = pgxscan.Get(ctx, r.db.Pool, &userID, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return userID, nil
}
