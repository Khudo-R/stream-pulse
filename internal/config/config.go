package config

import (
	"log"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type Config struct {
	AppPort     string `env:"APP_PORT" envDefault:"8080"`
	DatabaseURL string `env:"DATABASE_URL,required"`
	RabbitURL   string `env:"RABBIT_URL,required"`
	RedisURL    string `env:"REDIS_URL,required"`
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("Unable to parse config: %v", err)
	}

	return cfg
}
