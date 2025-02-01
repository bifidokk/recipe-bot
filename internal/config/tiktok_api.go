package config

import (
	"errors"
	"os"
)

const (
	// #nosec G101
	tikTokAPITokenEnvName = "TIK_TOK_TOKEN"
)

type tikTokAPIConfig struct {
	token string
}

func NewTikTokAPIConfig() (TikTokAPIConfig, error) {
	token := os.Getenv(tikTokAPITokenEnvName)
	if len(token) == 0 {
		return nil, errors.New("tik tok token not found")
	}

	return &tikTokAPIConfig{
		token: token,
	}, nil
}

func (cfg *tikTokAPIConfig) Token() string {
	return cfg.token
}
