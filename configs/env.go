package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DbName             string
	GithubClientId     string
	GithubClientSecret string
}

func LoadEnv() error {
	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}
	return nil
}

func GetConfig() (*Config, error) {
	err := LoadEnv()
	if err != nil {
		return nil, err
	}

	cfg := Config{
		DbName:             os.Getenv("DB_NAME"),
		GithubClientId:     os.Getenv("GITHUB_CLIENT_ID"),
		GithubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	}

	return &cfg, nil
}
