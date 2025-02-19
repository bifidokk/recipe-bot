package repository

import (
	"context"
	"fmt"

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
	tableName = "\"user\""

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

func (r *UserRepository) FindByTelegramID(ctx context.Context, id int) (*entity.User, error) {
	var user entity.User
	query, args, err := r.sqlBuilder.Select(
		idColumn, nameColumn, tgIDColumn, createdAtColumn, updatedAtColumn,
	).
		From(tableName).
		Where(sq.Eq{tgIDColumn: id}).
		ToSql()

	if err != nil {
		log.Error().Err(err).Msg("failed to build query")
		return nil, err
	}

	err = pgxscan.Select(ctx, r.db.Pool, &user, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user by ID: %w", err)
	}

	return &user, nil
}
