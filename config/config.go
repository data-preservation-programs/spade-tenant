package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

type TenantServiceConfig struct {
	PORT                string `env:"PORT" envDefault:"3000"`
	DB_URL              string `env:"DB_URL" envDefault:""`
	DRY_RUN             bool   `env:"DRY_RUN" envDefault:"false"`
	DEBUG               bool   `env:"DEBUG" envDefault:"false"`
	DB_ALLOW_MIGRATIONS bool   `env:"DB_ALLOW_MIGRATIONS" envDefault:"true"`
	BROKER_ACCESS_TOKEN string `env:"BROKER_ACCESS_TOKEN"`
}

func InitConfig() TenantServiceConfig {
	godotenv.Load() // load from environment OR .env file if it exists
	var cfg TenantServiceConfig

	if err := env.Parse(&cfg); err != nil {
		log.Fatal("error parsing config: %+v\n", err)
	}

	log.Debug("config parsed successfully")

	return cfg
}
