package service

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	kafka_lib "github.com/s21platform/kafka-lib"
	logger_lib "github.com/s21platform/logger-lib"
	optionhubproto "github.com/s21platform/optionhub-proto/optionhub-proto"
	optionhubproto_v1 "github.com/s21platform/optionhub-proto/optionhub/v1"

	"github.com/s21platform/optionhub-service/internal/config"
	"github.com/s21platform/optionhub-service/internal/model"
)

type Service struct {
	optionhubproto.UnimplementedOptionhubServiceServer
	dbR      DBRepo
	setAttrP *kafka_lib.KafkaProducer
}

func NewService(repo DBRepo, setAttributeProducer *kafka_lib.KafkaProducer) *Service {
	return &Service{dbR: repo, setAttrP: setAttributeProducer}
}

func (s *Service) GetOsByID(ctx context.Context, in *optionhubproto.GetByIdIn) (*optionhubproto.GetByIdOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetOsByID")

	os, err := s.dbR.GetOsByID(ctx, in.Id)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get os by id, err: %v", err))
		return nil, status.Errorf(codes.NotFound, "failed to get os by id, err: %v", err)
	}

	return &optionhubproto.GetByIdOut{Id: in.Id, Value: os}, nil
}

func (s *Service) AddOs(ctx context.Context, in *optionhubproto.AddIn) (*optionhubproto.AddOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("AddOs")

	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		logger.Error("failed to find uuid")
		return nil, status.Errorf(codes.Unauthenticated, "failed to find uuid")
	}

	id, err := s.dbR.AddOS(ctx, in.Value, uuid)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to add new os, err: %v", err))
		return nil, status.Errorf(codes.Aborted, "failed to add new os, err: %v", err)
	}

	return &optionhubproto.AddOut{Id: id, Value: in.Value}, nil
}

func (s *Service) GetOsBySearchName(ctx context.Context, in *optionhubproto.GetByNameIn) (*optionhubproto.GetByNameOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetOsBySearchName")

	var (
		osList model.CategoryItemList
		err    error
	)

	if len(in.Name) < 2 {
		osList, err = s.dbR.GetOsPreview(ctx)
	} else {
		osList, err = s.dbR.GetOsBySearchName(ctx, in.Name)
	}

	if err != nil {
		logger.Error(fmt.Sprintf("failed to get os by name, err: %v", err))
		return nil, status.Errorf(codes.NotFound, "failed to get os by name, err: %v", err)
	}

	return &optionhubproto.GetByNameOut{
		Options: osList.FromDTO(),
	}, nil
}

func (s *Service) GetAllOs(ctx context.Context, _ *optionhubproto.EmptyOptionhub) (*optionhubproto.GetAllOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetAllOs")

	OSList, err := s.dbR.GetAllOs()
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get all os list: %v", err))
		return nil, status.Errorf(codes.NotFound, "failed to get all os list: %v", err)
	}

	return &optionhubproto.GetAllOut{
		Values: OSList.FromDTO(),
	}, nil
}

func (s *Service) SetAttribute(ctx context.Context, in *optionhubproto_v1.SetAttributeByIdIn) error {
  logger := logger_lib.FromContext(ctx, config.KeyLogger)
  logger.AddFuncName("SetAttribute")

  err := s.dbR.SetAttribute(ctx, in)
  if err != nil {
    logger.Error(fmt.Sprintf("failed to add new attribute: %v", err))
    return status.Errorf(codes.Aborted, "failed to add new attribute: %v", err)
  }

  message := model.SetAttributeMessage{AttributeId: in.AttributeId}
  err = s.setAttrP.ProduceMessage(message)
  if err != nil {
    logger.Error(fmt.Sprintf("failed to produce kafka message: %v", err))
    return status.Errorf(codes.Aborted, "failed to produce kafka message: %v", err)
  }

  return nil
}