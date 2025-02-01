package config

import (
	"errors"
	"os"
)

const (
	// #nosec G101
	tokenEnvName = "TG_BOT_TOKEN"
)

type tgBotConfig struct {
	token string
}

func NewTgBotConfig() (TgBotConfig, error) {
	token := os.Getenv(tokenEnvName)
	if len(token) == 0 {
		return nil, errors.New("tg bot token not found")
	}

	return &tgBotConfig{
		token: token,
	}, nil
}

func (cfg *tgBotConfig) Token() string {
	return cfg.token
}
