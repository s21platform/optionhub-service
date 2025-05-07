package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Service  Service
	Postgres Postgres
	Metrics  Metrics
	Logger   Logger
	Platform Platform
	Kafka    Kafka
}

type Service struct {
	Port string `env:"OPTIONHUB_SERVICE_PORT"`
	Name string `env:"OPTIONHUB_SERVICE_NAME"`
}

type Postgres struct {
	User     string `env:"OPTIONHUB_SERVICE_POSTGRES_USER"`
	Password string `env:"OPTIONHUB_SERVICE_POSTGRES_PASSWORD"`
	Database string `env:"OPTIONHUB_SERVICE_POSTGRES_DB"`
	Host     string `env:"OPTIONHUB_SERVICE_POSTGRES_HOST"`
	Port     string `env:"OPTIONHUB_SERVICE_POSTGRES_PORT"`
}

type Metrics struct {
	Host string `env:"GRAFANA_HOST"`
	Port int    `env:"GRAFANA_PORT"`
}

type Logger struct {
	Host string `env:"LOGGER_SERVICE_HOST"`
	Port string `env:"LOGGER_SERVICE_PORT"`
}

type Platform struct {
	Env string `env:"ENV"` // окружение (stage)
}

type Kafka struct {
	Host              string `env:"KAFKA_HOST"`
	Port              string `env:"KAFKA_PORT"`
	SetAttributeTopic string `env:"STAFF_SET_ATTRIBUTE"`
}

func NewConfig() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)

	if err != nil {
		log.Fatalf("failed to read env variables: %s", err)
	}

	return cfg
}
