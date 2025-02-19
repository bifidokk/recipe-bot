package client

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type DBClient struct {
	Pool *pgxpool.Pool
}

func NewDBClient(dsn string) (*DBClient, error) {
	config, err := pgxpool.ParseConfig(dsn)

	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnIdleTime = 5 * time.Minute
	config.MaxConnLifetime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	log.Info().Msg("Database connection established successfully")

	return &DBClient{Pool: pool}, nil
}

func (db *DBClient) Close() {
	db.Pool.Close()
	log.Info().Msg("Database connection closed")
}
