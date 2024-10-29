package config

import (
	"errors"
	"os"
)

const (
	openAIApiTokenEnvName = "OPEN_AI_TOKEN"
)

type openAIApiConfig struct {
	token string
}

func NewOpenAIAPIConfig() (OpenAIAPIConfig, error) {
	token := os.Getenv(openAIApiTokenEnvName)
	if len(token) == 0 {
		return nil, errors.New("open ai token not found")
	}

	return &openAIApiConfig{
		token: token,
	}, nil
}

func (cfg *openAIApiConfig) Token() string {
	return cfg.token
}
