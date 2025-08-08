package config

import (
	"github.com/joho/godotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"log"
	"strings"
)

var k = koanf.New(".")

func LoadConfig() (*Config, error) {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Println("No .env file found, using system env vars instead")
	}

	provider := env.Provider("", ".", func(s string) string {
		return strings.ToLower(s)
	})

	err = k.Load(provider, nil)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = k.Unmarshal("", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
