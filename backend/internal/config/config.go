package config

import "github.com/caarlos0/env/v11"

type Config struct {
	Port        int    `env:"PORT" envDefault:"8000"`
	DatabaseURL string `env:"DATABASE_URL,required"`
}

func Load() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
