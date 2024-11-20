package service

import (
	"context"
	"fmt"
	"optionhub-service/internal/config"

	optionhubproto "github.com/s21platform/optionhub-proto/optionhub-proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	// optionhubproto.OptionhubServiceClient
	optionhubproto.UnimplementedOptionhubServiceServer
	dbR DBRepo
}

func NewService(repo DBRepo) *Service {
	return &Service{dbR: repo}
}

func (s *Service) GetOsByID(ctx context.Context, in *optionhubproto.GetByIdIn) (*optionhubproto.GetByIdOut, error) {
	os, err := s.dbR.GetOsByID(ctx, in.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cannot get os by id, err: %v", err)
	}

	return &optionhubproto.GetByIdOut{Id: in.Id, Value: os}, nil
}

func (s *Service) AddOs(ctx context.Context, in *optionhubproto.AddIn) (*optionhubproto.AddOut, error) {
	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "cannot find uuid")
	}

	id, err := s.dbR.AddOS(ctx, in.Value, uuid)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "cannot add new os, err: %v", err)
	}

	return &optionhubproto.AddOut{Id: id, Value: in.Value}, nil
}

func (s *Service) GetOsBySearchName(ctx context.Context, in *optionhubproto.GetByNameIn) (*optionhubproto.GetByNameOut,
	error) {
	OS, err := s.dbR.GetOsBySearchName(ctx, in.Name)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cannot get os by name, err: %v", err)
	}

	return &optionhubproto.GetByNameOut{Values: OS}, nil
}

func (s *Service) GetAllOs(ctx context.Context, in *optionhubproto.GetAllIn) (*optionhubproto.GetAllOut, error) {
	_, _ = ctx, in

	OSList, err := s.dbR.GetAllOs()
	if err != nil {
		return nil, fmt.Errorf("failed to get all os list: %w", err)
	}

	return &optionhubproto.GetAllOut{Values: OSList}, nil
}
