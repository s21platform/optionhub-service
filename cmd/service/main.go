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
	"optionhub-service/internal/repository/postgres"

	optionhubproto "github.com/s21platform/optionhub-proto/optionhub-proto"
)

func main() {
	cfg := config.NewConfig()

	dbRepo := postgres.New(cfg)
	defer dbRepo.Close()

	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, "optionhub", cfg.Platform.Env)
	if err != nil {
		log.Fatalf("failed to create metrics: %v", err)
	}
	defer metrics.Disconnect()

	optionhubService := service.NewService(dbRepo)

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			infra.AuthInterceptor,
			infra.MetricsInterceptor(metrics),
		),
	)
	optionhubproto.RegisterOptionhubServiceServer(s, optionhubService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Printf("failed to listen port: %s; Error: %s", cfg.Service.Port, err)
	}

	if err = s.Serve(lis); err != nil {
		log.Printf("failed to start service: %s; Error: %s", cfg.Service.Port, err)
	}
}
