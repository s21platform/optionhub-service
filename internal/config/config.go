package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Service  Service
	Postgres Postgres
}

type Service struct {
	Port string `env:"OPTIONHUB_SERVICE_PORT"`
}

type Postgres struct {
	User     string `env:"OPTIONHUB_SERVICE_POSTGRES_USER"`
	Password string `env:"OPTIONHUB_SERVICE_POSTGRES_PASSWORD"`
	Database string `env:"OPTIONHUB_SERVICE_POSTGRES_DB"`
	Host     string `env:"OPTIONHUB_SERVICE_POSTGRES_HOST"`
	Port     string `env:"OPTIONHUB_SERVICE_POSTGRES_PORT"`
}

func NewConfig() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatalf("Ну пиздец ошибка, %s", err.Error())
	}
	return cfg
}
