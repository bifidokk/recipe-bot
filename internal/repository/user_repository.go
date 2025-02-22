package repository

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

type UserRepository struct {
	db         *client.DBClient
	sqlBuilder sq.StatementBuilderType
}

const (
	tableName = "users"

	idColumn        = "id"
	nameColumn      = "name"
	tgIDColumn      = "telegram_id"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

func NewUserRepository(dbClient *client.DBClient) *UserRepository {
	return &UserRepository{
		db:         dbClient,
		sqlBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *UserRepository) FindByTelegramID(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User
	query, args, err := r.sqlBuilder.Select(
		idColumn, nameColumn, tgIDColumn, createdAtColumn, updatedAtColumn,
	).
		From(tableName).
		Where(sq.Eq{tgIDColumn: strconv.Itoa(int(id))}).
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

func (r *UserRepository) FindByID(ctx context.Context, id int) (*entity.User, error) {
	var user entity.User
	query, args, err := r.sqlBuilder.Select(
		idColumn, nameColumn, tgIDColumn, createdAtColumn, updatedAtColumn,
	).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
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

func (r *UserRepository) CreateUser(ctx context.Context, user *entity.User) (int, error) {
	query, args, err := r.sqlBuilder.Insert(tableName).
		Columns(nameColumn, tgIDColumn).
		Values(user.Name, user.TelegramID).
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
