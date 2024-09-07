package main

import (
	"errors"
	"log"

	"optionhub-service/internal/config"
	"optionhub-service/internal/repository/db"
)

func main() {
	cfg := config.NewConfig()

	dbRepo, err := db.New(cfg)
	if err != nil {
		log.Fatal(errors.New("нет подключения к бд"))
	}

	defer dbRepo.Close()
}
