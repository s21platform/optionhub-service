package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Service Service
}

type Service struct {
	Port string `env:"OPTIONHUB_SERVICE_PORT"`
	Host string `env:"OPTIONHUB_SERVICE_HOST"`
}

func NewConfig() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatalf("Ну пиздец ошибка, %s", err.Error())
	}
	return cfg
}
