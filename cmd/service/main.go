package main

import (
	"fmt"
	"log"
	"net"
	"optionhub-service/internal/infra"
	"optionhub-service/internal/service"

	"google.golang.org/grpc"

	"github.com/s21platform/metrics-lib/pkg"

	"optionhub-service/internal/config"
	"optionhub-service/internal/repository/db"

	optionhubproto "github.com/s21platform/optionhub-proto/optionhub-proto"
)

func main() {
	cfg := config.NewConfig()

	dbRepo, err := db.New(cfg)
	if err != nil {
		log.Fatalf("Error initialize db repository: %v", err)
	}
	defer dbRepo.Close()

	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, "optionhub", cfg.Platform.Env)
	if err != nil {
		log.Fatalf("fail to create metrics: %v", err)
	}
	defer metrics.Disconnect()

	optionhubService := service.NewService(dbRepo)

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			infra.UnaryInterceptor,
			infra.MetricsInterceptor(metrics),
		),
	)
	optionhubproto.RegisterOptionhubServiceServer(s, optionhubService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Printf("Cannot listen port: %s; Error: %s", cfg.Service.Port, err)
	}

	if err = s.Serve(lis); err != nil {
		log.Printf("Cannot start service: %s; Error: %s", cfg.Service.Port, err)
	}
}
