package service

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	logger_lib "github.com/s21platform/logger-lib"
	optionhubproto "github.com/s21platform/optionhub-proto/optionhub-proto"

	"optionhub-service/internal/config"
	"optionhub-service/internal/model"
)

type Service struct {
	optionhubproto.UnimplementedOptionhubServiceServer
	dbR DBRepo
}

func NewService(repo DBRepo) *Service {
	return &Service{dbR: repo}
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
