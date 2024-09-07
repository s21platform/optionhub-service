package main

import (
	"log"

	"optionhub-service/internal/config"
	"optionhub-service/internal/repository/db"
)

func main() {
	cfg := config.NewConfig()

	dbRepo, err := db.New(cfg)
	if err != nil {
		log.Fatalf("Error initialize db repository: %v", err)
	}

	defer dbRepo.Close()
}
