package service

import (
	"context"
	"optionhub-service/internal/config"

	optionhubproto "github.com/s21platform/optionhub-proto/optionhub-proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	//optionhubproto.OptionhubServiceClient
	optionhubproto.UnimplementedOptionhubServiceServer
	dbR DbRepo
}

func NewService(repo DbRepo) *Service {
	return &Service{dbR: repo}
}

func (s *Service) GetOsById(ctx context.Context, in *optionhubproto.GetByIdIn) (*optionhubproto.GetByIdOut, error) {
	os, err := s.dbR.GetOsById(ctx, in.Id)
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

func (s *Service) GetOsBySearchName(ctx context.Context, in *optionhubproto.GetByNameIn) (*optionhubproto.GetByNameOut, error) {
	if len(in.Name) < 2 {
		return nil, nil
	}

	os, err := s.dbR.GetOsBySearchName(ctx, in.Name)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cannot get os by name, err: %v", err)
	}

	var records optionhubproto.GetByNameOut
	for _, osObj := range os {
		records.Values = append(records.Values,
			&optionhubproto.Record{
				Id:    osObj.Id,
				Value: osObj.Name,
			})

	}

	return &records, nil
}
