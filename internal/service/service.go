package service

import (
	"context"
	optionhub_proto "github.com/s21platform/optionhub-proto/optionhub-proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	optionhub_proto.OptionhubServiceClient
	dbR DbRepo
}

func NewService(repo DbRepo) *Service {
	return &Service{dbR: repo}
}

func (s *Service) GetOsById(ctx context.Context, in *optionhub_proto.GetByIdIn) (*optionhub_proto.GetByIdOut, error) {
	os, err := s.dbR.GetOsById(ctx, in.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cannot get os by id, err: %v", err)
	}
	return &optionhub_proto.GetByIdOut{Id: in.Id, Value: os}, nil
}

func (s *Service) AddOs(ctx context.Context, in *optionhub_proto.AddIn) (*optionhub_proto.AddOut, error) {
	id, err := s.dbR.AddOS(ctx, in.Value)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "cannot add new os, err: %v", err)

	}
	return &optionhub_proto.AddOut{Id: id, Value: in.Value}, nil
}

func (s *Service) GetAllOs(ctx context.Context, in *optionhub_proto.AddIn) (*optionhub_proto.AddOut, error) {
	return nil, nil
}
