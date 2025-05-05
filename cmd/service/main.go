package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	kafka_lib "github.com/s21platform/kafka-lib"
	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/metrics-lib/pkg"
	optionhubproto "github.com/s21platform/optionhub-proto/optionhub-proto"
	optionhubprotov1 "github.com/s21platform/optionhub-proto/optionhub/v1"

	"github.com/s21platform/optionhub-service/internal/config"
	"github.com/s21platform/optionhub-service/internal/infra"
	"github.com/s21platform/optionhub-service/internal/repository/postgres"
	"github.com/s21platform/optionhub-service/internal/service"
	servicev1 "github.com/s21platform/optionhub-service/internal/service/v1"
)

func main() {
	cfg := config.NewConfig()
	logger := logger_lib.New(cfg.Logger.Host, cfg.Logger.Port, cfg.Service.Name, cfg.Platform.Env)

	dbRepo := postgres.New(cfg)
	defer dbRepo.Close()

	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, "optionhub", cfg.Platform.Env)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to create metrics: %v", err))
		log.Fatalf("failed to create metrics: %v", err)
	}
	defer metrics.Disconnect()

	kafkaConfig := kafka_lib.DefaultProducerConfig(cfg.Kafka.Host, cfg.Kafka.Port, cfg.Kafka.SetAttributeTopic)

	producerSetAttribute := kafka_lib.NewProducer(kafkaConfig)
	defer func(producerSetAttribute *kafka_lib.KafkaProducer) {
		err := producerSetAttribute.Close()
		if err != nil {
			logger.Error(fmt.Sprintf("failed to close producer: %v", err))
		}
	}(producerSetAttribute)

	optionhubService := service.NewService(dbRepo)
	optionhubServicev1 := servicev1.NewService(dbRepo, producerSetAttribute)

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			infra.AuthInterceptor,
			infra.MetricsInterceptor(metrics),
			infra.Logger(logger),
		),
	)
	optionhubproto.RegisterOptionhubServiceServer(s, optionhubService)
	optionhubprotov1.RegisterOptionhubServiceV1Server(s, optionhubServicev1)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to listen port: %s; Error: %s", cfg.Service.Port, err))
	}

	if err = s.Serve(lis); err != nil {
		logger.Error(fmt.Sprintf("failed to start service: %s; Error: %s", cfg.Service.Port, err))
	}
}
