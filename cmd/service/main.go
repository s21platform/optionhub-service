package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/go-chi/chi/v5"
	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/metrics-lib/pkg"
	"github.com/soheilhy/cmux"

	"github.com/s21platform/optionhub-service/internal/config"
	api "github.com/s21platform/optionhub-service/internal/generated"
	"github.com/s21platform/optionhub-service/internal/infra"
	"github.com/s21platform/optionhub-service/internal/repository/postgres"
	"github.com/s21platform/optionhub-service/internal/rest"
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

	grpcSrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			infra.AuthInterceptor,
			infra.MetricsInterceptor(metrics),
			infra.Logger(logger),
		),
	)
	optionhub.RegisterOptionhubServiceServer(grpcSrv, optionhubService)

	handler := rest.New(dbRepo)
	router := chi.NewRouter()
	router.Use(infra.AuthRequest)
	router.Use(infra.LoggerRequest(logger))
	router.Use(infra.MetricsRequest(metrics))

	api.HandlerFromMux(handler, router)
	httpServer := &http.Server{
		Handler: router,
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		logger_lib.Error(ctx, fmt.Sprintf("failed to listen port: %v", err))
	}

	m := cmux.New(lis)

	grpcListener := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpListener := m.Match(cmux.HTTP1Fast())

	g, _ := errgroup.WithContext(context.Background())

	logger_lib.Info(ctx, "starting server")
	g.Go(func() error {
		if err := grpcSrv.Serve(grpcListener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			return fmt.Errorf("gRPC server error: %v", err)
		}
		return nil
	})

	g.Go(func() error {
		if err := httpServer.Serve(httpListener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("HTTP server error: %v", err)
		}
		return nil
	})

	g.Go(func() error {
		if err := m.Serve(); err != nil {
			return fmt.Errorf("cannot start service: %v", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		logger_lib.Error(ctx, fmt.Sprintf("server error: %v", err))
		log.Fatalf("Server error: %v", err)
	}
}
