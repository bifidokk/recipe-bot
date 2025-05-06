package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

// Load is function to load env variables to os
func Load(path string) error {
	_ = godotenv.Load(path)

	requiredVars := []string{
		"TG_BOT_TOKEN",
		"TIK_TOK_TOKEN",
		"OPEN_AI_TOKEN",
		"PG_DSN",
	}

	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			return fmt.Errorf("required environment variable %s is not set", v)
		}
	}

	return nil
}
