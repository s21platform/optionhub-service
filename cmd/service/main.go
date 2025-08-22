package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/metrics-lib/pkg"

	"github.com/s21platform/optionhub-service/internal/config"
	"github.com/s21platform/optionhub-service/internal/infra"
	"github.com/s21platform/optionhub-service/internal/repository/postgres"
	"github.com/s21platform/optionhub-service/internal/service"
	"github.com/s21platform/optionhub-service/pkg/optionhub"
)

func main() {
	cfg := config.NewConfig()
	log.Println(cfg)
	logger := logger_lib.New(cfg.Logger.Host, cfg.Logger.Port, cfg.Service.Name, cfg.Platform.Env)
	ctx := logger_lib.NewContext(context.Background(), logger)

	dbRepo := postgres.New(cfg)
	defer dbRepo.Close()

	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, "optionhub", cfg.Platform.Env)
	if err != nil {
		logger_lib.Error(ctx, fmt.Sprintf("failed to create metrics: %v", err))
		log.Fatalf("failed to create metrics: %v", err)
	}
	defer metrics.Disconnect()

	optionhubService := service.NewService(dbRepo)

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			infra.AuthInterceptor,
			infra.MetricsInterceptor(metrics),
			infra.Logger(logger),
		),
	)

	optionhub.RegisterOptionhubServiceServer(s, optionhubService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		logger_lib.Error(ctx, fmt.Sprintf("failed to listen port: %s; Error: %s", cfg.Service.Port, err))
	}

	if err = s.Serve(lis); err != nil {
		logger_lib.Error(ctx, fmt.Sprintf("failed to start service: %s; Error: %s", cfg.Service.Port, err))
	}
}
