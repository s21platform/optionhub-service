package service

import (
	"context"
	optionhubproto "github.com/s21platform/optionhub-proto/optionhub-proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	//optionhubproto.OptionhubServiceClient
	optionhubproto.UnimplementedOptionhubServiceServer
	dbR DbRepo
}

func (s *Service) GetAllOs(ctx context.Context, in *optionhubproto.GetAllIn) (*optionhubproto.GetAllOut, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) SetOsById(ctx context.Context, in *optionhubproto.SetByIdIn) (*optionhubproto.SetByIdOut, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) DeleteOsById(ctx context.Context, in *optionhubproto.DeleteByIdIn) (*optionhubproto.DeleteByIdOut, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetWorkPlaceById(ctx context.Context, in *optionhubproto.GetByIdIn) (*optionhubproto.GetByIdOut, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetAllWorkPlace(ctx context.Context, in *optionhubproto.GetAllIn) (*optionhubproto.GetAllOut, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) AddWorkPlace(ctx context.Context, in *optionhubproto.AddIn) (*optionhubproto.AddOut, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) SetWorkPlaceById(ctx context.Context, in *optionhubproto.SetByIdIn) (*optionhubproto.SetByIdOut, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) DeleteWorkPlaceById(ctx context.Context, in *optionhubproto.DeleteByIdIn) (*optionhubproto.DeleteByIdOut, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetStudyPlaceById(ctx context.Context, in *optionhubproto.GetByIdIn) (*optionhubproto.GetByIdOut, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetAllStudyPlace(ctx context.Context, in *optionhubproto.GetAllIn) (*optionhubproto.GetAllOut, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) AddStudyPlace(ctx context.Context, in *optionhubproto.AddIn) (*optionhubproto.AddOut, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) SetStudyPlaceById(ctx context.Context, in *optionhubproto.SetByIdIn) (*optionhubproto.SetByIdOut, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) DeleteStudyPlaceById(ctx context.Context, in *optionhubproto.DeleteByIdIn) (*optionhubproto.DeleteByIdOut, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetHobbyById(ctx context.Context, in *optionhubproto.GetByIdIn) (*optionhubproto.GetByIdOut, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetHobbyPlace(ctx context.Context, in *optionhubproto.GetAllIn) (*optionhubproto.GetAllOut, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) AddHobby(ctx context.Context, in *optionhubproto.AddIn) (*optionhubproto.AddOut, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) SetHobbyById(ctx context.Context, in *optionhubproto.SetByIdIn) (*optionhubproto.SetByIdOut, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) DeleteHobbyById(ctx context.Context, in *optionhubproto.DeleteByIdIn) (*optionhubproto.DeleteByIdOut, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetSkillById(ctx context.Context, in *optionhubproto.GetByIdIn) (*optionhubproto.GetByIdOut, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetSkillPlace(ctx context.Context, in *optionhubproto.GetAllIn) (*optionhubproto.GetAllOut, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) AddSkill(ctx context.Context, in *optionhubproto.AddIn) (*optionhubproto.AddOut, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) SetSkillById(ctx context.Context, in *optionhubproto.SetByIdIn) (*optionhubproto.SetByIdOut, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) DeleteSkillById(ctx context.Context, in *optionhubproto.DeleteByIdIn) (*optionhubproto.DeleteByIdOut, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetCityById(ctx context.Context, in *optionhubproto.GetByIdIn) (*optionhubproto.GetByIdOut, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetCityPlace(ctx context.Context, in *optionhubproto.GetAllIn) (*optionhubproto.GetAllOut, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) AddCity(ctx context.Context, in *optionhubproto.AddIn) (*optionhubproto.AddOut, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) SetCityById(ctx context.Context, in *optionhubproto.SetByIdIn) (*optionhubproto.SetByIdOut, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) DeleteCityById(ctx context.Context, in *optionhubproto.DeleteByIdIn) (*optionhubproto.DeleteByIdOut, error) {
	//TODO implement me
	panic("implement me")
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
	id, err := s.dbR.AddOS(ctx, in.Value)
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
	for idx := range os {
		records.Values = append(records.Values,
			&optionhubproto.Record{
				Id:    os[idx].Id,
				Value: os[idx].Name,
			})

	}

	return &records, nil
}
