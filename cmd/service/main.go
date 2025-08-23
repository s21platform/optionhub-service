package main

import (
	"context"
	"fmt"
	"net"
	"os"

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
	logger := logger_lib.New(cfg.Logger.Host, cfg.Logger.Port, cfg.Service.Name, cfg.Platform.Env)
	ctx := logger_lib.NewContext(context.Background(), logger)

	dbRepo := postgres.New(cfg)
	defer dbRepo.Close()

	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, cfg.Service.Name, cfg.Platform.Env)
	if err != nil {
		logger_lib.Error(ctx, fmt.Sprintf("failed to create metrics: %v", err))
		os.Exit(1)
	}
	defer metrics.Disconnect()

	optionhubService := service.NewService(dbRepo)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			infra.AuthInterceptor,
			infra.MetricsInterceptor(metrics),
			infra.Logger(logger),
		),
	)

	optionhub.RegisterOptionhubServiceServer(server, optionhubService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		logger_lib.Error(ctx, fmt.Sprintf("failed to listen port: %v", err))
	}

	if err = server.Serve(lis); err != nil {
		logger_lib.Error(ctx, fmt.Sprintf("failed to start service: %v", err))
	}
}
