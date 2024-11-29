package service

import (
	"context"
	"optionhub-service/internal/config"
	"optionhub-service/internal/model"

	optionhubproto "github.com/s21platform/optionhub-proto/optionhub-proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
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
		return nil, status.Errorf(codes.NotFound, "cannot get os by name, err: %v", err)
	}

	return &optionhubproto.GetByNameOut{
		Options: osList.FromDTO(),
	}, nil
}

func (s *Service) GetAllOs(ctx context.Context, in *optionhubproto.EmptyOptionhub) (*optionhubproto.GetAllOut, error) {
	_, _ = ctx, in

	OSList, err := s.dbR.GetAllOs()
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get all os list: %v", err)
	}

	return &optionhubproto.GetAllOut{
		Values: OSList.FromDTO(),
	}, nil
}

func (s *Service) GetWorkPlaceBySearchName(ctx context.Context, in *optionhubproto.GetByNameIn) (
	*optionhubproto.GetByNameOut, error) {
	var (
		workPlaceList model.CategoryItemList
		err           error
	)

	if len(in.Name) < 2 {
		workPlaceList, err = s.dbR.GetWorkPlacePreview(ctx)
	} else {
		workPlaceList, err = s.dbR.GetWorkPlaceBySearchName(ctx, in.Name)
	}

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get workplace list: %v", err)
	}

	return &optionhubproto.GetByNameOut{
		Options: workPlaceList.FromDTO(),
	}, nil
}

func (s *Service) GetWorkPlaceById(ctx context.Context, in *optionhubproto.GetByIdIn) (
	*optionhubproto.GetByIdOut, error) {
	workplace, err := s.dbR.GetWorkPlaceByID(ctx, in.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get workplace by id: %v", err)
	}

	return &optionhubproto.GetByIdOut{Id: in.Id, Value: workplace}, nil
}

func (s *Service) AddWorkPlace(ctx context.Context, in *optionhubproto.AddIn) (*optionhubproto.AddOut, error) {
	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "failed to find uuid")
	}

	id, err := s.dbR.AddWorkPlace(ctx, in.Value, uuid)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "failed to add new workplace: %v", err)
	}

	return &optionhubproto.AddOut{Id: id, Value: in.Value}, nil
}

func (s *Service) GetStudyPlaceBySearchName(ctx context.Context, in *optionhubproto.GetByNameIn) (
	*optionhubproto.GetByNameOut, error) {
	var (
		studyPlaceList model.CategoryItemList
		err            error
	)

	if len(in.Name) < 2 {
		studyPlaceList, err = s.dbR.GetStudyPlacePreview(ctx)
	} else {
		studyPlaceList, err = s.dbR.GetStudyPlaceBySearchName(ctx, in.Name)
	}

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get study place list: %v", err)
	}

	return &optionhubproto.GetByNameOut{
		Options: studyPlaceList.FromDTO(),
	}, nil
}

func (s *Service) GetStudyPlaceById(ctx context.Context, in *optionhubproto.GetByIdIn) (
	*optionhubproto.GetByIdOut, error) {
	studyPlace, err := s.dbR.GetStudyPlaceByID(ctx, in.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get study place by id: %v", err)
	}

	return &optionhubproto.GetByIdOut{Id: in.Id, Value: studyPlace}, nil
}

func (s *Service) AddStudyPlace(ctx context.Context, in *optionhubproto.AddIn) (*optionhubproto.AddOut, error) {
	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "failed to find uuid")
	}

	id, err := s.dbR.AddStudyPlace(ctx, in.Value, uuid)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "failed to add new study place: %v", err)
	}

	return &optionhubproto.AddOut{Id: id, Value: in.Value}, nil
}

func (s *Service) GetHobbyBySearchName(ctx context.Context, in *optionhubproto.GetByNameIn) (
	*optionhubproto.GetByNameOut, error) {
	var (
		hobbyList model.CategoryItemList
		err       error
	)

	if len(in.Name) < 2 {
		hobbyList, err = s.dbR.GetHobbyPreview(ctx)
	} else {
		hobbyList, err = s.dbR.GetHobbyBySearchName(ctx, in.Name)
	}

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get hobby list: %v", err)
	}

	return &optionhubproto.GetByNameOut{
		Options: hobbyList.FromDTO(),
	}, nil
}

func (s *Service) GetHobbyById(ctx context.Context, in *optionhubproto.GetByIdIn) (*optionhubproto.GetByIdOut, error) {
	hobby, err := s.dbR.GetHobbyByID(ctx, in.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get hobby by id: %v", err)
	}

	return &optionhubproto.GetByIdOut{Id: in.Id, Value: hobby}, nil
}

func (s *Service) AddHobby(ctx context.Context, in *optionhubproto.AddIn) (*optionhubproto.AddOut, error) {
	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "failed to find uuid")
	}

	id, err := s.dbR.AddHobby(ctx, in.Value, uuid)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "failed to add new hobby: %v", err)
	}

	return &optionhubproto.AddOut{Id: id, Value: in.Value}, nil
}

func (s *Service) GetSkillBySearchName(ctx context.Context, in *optionhubproto.GetByNameIn) (
	*optionhubproto.GetByNameOut, error) {
	var (
		skillList model.CategoryItemList
		err       error
	)

	if len(in.Name) < 2 {
		skillList, err = s.dbR.GetSkillPreview(ctx)
	} else {
		skillList, err = s.dbR.GetSkillBySearchName(ctx, in.Name)
	}

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get skill list: %v", err)
	}

	return &optionhubproto.GetByNameOut{
		Options: skillList.FromDTO(),
	}, nil
}

func (s *Service) GetSkillById(ctx context.Context, in *optionhubproto.GetByIdIn) (*optionhubproto.GetByIdOut, error) {
	skill, err := s.dbR.GetSkillByID(ctx, in.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get skill by id: %v", err)
	}

	return &optionhubproto.GetByIdOut{Id: in.Id, Value: skill}, nil
}

func (s *Service) AddSkill(ctx context.Context, in *optionhubproto.AddIn) (*optionhubproto.AddOut, error) {
	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "failed to find uuid")
	}

	id, err := s.dbR.AddSkill(ctx, in.Value, uuid)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "failed to add new skill: %v", err)
	}

	return &optionhubproto.AddOut{Id: id, Value: in.Value}, nil
}

func (s *Service) GetCityBySearchName(ctx context.Context, in *optionhubproto.GetByNameIn) (
	*optionhubproto.GetByNameOut, error) {
	var (
		cityList model.CategoryItemList
		err      error
	)

	if len(in.Name) < 2 {
		cityList, err = s.dbR.GetCityPreview(ctx)
	} else {
		cityList, err = s.dbR.GetCityBySearchName(ctx, in.Name)
	}

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get city list: %v", err)
	}

	return &optionhubproto.GetByNameOut{
		Options: cityList.FromDTO(),
	}, nil
}

func (s *Service) GetCityById(ctx context.Context, in *optionhubproto.GetByIdIn) (*optionhubproto.GetByIdOut, error) {
	city, err := s.dbR.GetCityByID(ctx, in.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get city by id: %v", err)
	}

	return &optionhubproto.GetByIdOut{Id: in.Id, Value: city}, nil
}

func (s *Service) AddCity(ctx context.Context, in *optionhubproto.AddIn) (*optionhubproto.AddOut, error) {
	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "failed to find uuid")
	}

	id, err := s.dbR.AddCity(ctx, in.Value, uuid)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "failed to add new city: %v", err)
	}

	return &optionhubproto.AddOut{Id: id, Value: in.Value}, nil
}
