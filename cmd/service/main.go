package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"optionhub-service/internal/service"

	optionhubproto "github.com/s21platform/optionhub-proto/optionhub-proto"
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

	optionhubService := service.NewService(dbRepo)

	s := grpc.NewServer()
	optionhubproto.RegisterOptionhubServiceServer(s, optionhubService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Printf("Cannot listen port: %s; Error: %s", cfg.Service.Port, err)
	}

	if err = s.Serve(lis); err != nil {
		log.Printf("Cannot start service: %s; Error: %s", cfg.Service.Port, err)
	}
}
