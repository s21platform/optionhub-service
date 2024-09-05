package service

import (
	"context"
	"fmt"
	optionhub_proto "github.com/s21platform/optionhub-proto/optionhub-proto"
	"optionhub-service/internal/repository/db"
)

type service struct {
	optionhub_proto.OptionhubServiceClient
	db db.DbRepo
}

func NewService(repo db.DbRepo) *service {
	return &service{db: repo}
}

func (s *service) GetOsById(ctx context.Context, in *optionhub_proto.GetByIdIn) (*optionhub_proto.GetByIdOut, error) {
	os, err := s.db.GetOsById(ctx, in.Id)
	if err != nil {
		return nil, fmt.Errorf("cannot get os by id, err: %v", err)
	}
	return &optionhub_proto.GetByIdOut{Id: in.Id, Value: os}, nil
}

func (s *service) AddOs(ctx context.Context, in *optionhub_proto.AddIn) (*optionhub_proto.AddOut, error) {
	id, err := s.db.AddOS(ctx, in.Value)
	if err != nil {
		return nil, fmt.Errorf("cannot add new os, err: %v", err)
	}
	return &optionhub_proto.AddOut{Id: id, Value: in.Value}, nil
}

func (s *service) GetAllOs(ctx context.Context, in *optionhub_proto.AddIn) (*optionhub_proto.AddOut, error) {
	return nil, nil
}
