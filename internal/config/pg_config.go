package config

import (
	"errors"
	"os"
)

const (
	pgDsnEnvName = "PG_DSN"
)

type pgConfig struct {
	dsn string
}

func NewPgConfig() (PgConfig, error) {
	dsn := os.Getenv(pgDsnEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn is not found")
	}

	return &pgConfig{
		dsn: dsn,
	}, nil
}

func (c *pgConfig) Dsn() string {
	return c.dsn
}
