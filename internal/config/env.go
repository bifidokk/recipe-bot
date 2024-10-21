package config

import (
	"github.com/joho/godotenv"
)

// Load is function to load env variables to os
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}
