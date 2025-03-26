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

	"optionhub-service/internal/config"
	"optionhub-service/internal/model"
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

func (s *Service) GetWorkPlaceBySearchName(ctx context.Context, in *optionhubproto.GetByNameIn) (*optionhubproto.GetByNameOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetWorkPlaceBySearchName")

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
		logger.Error(fmt.Sprintf("failed to get workplace list: %v", err))
		return nil, status.Errorf(codes.NotFound, "failed to get workplace list: %v", err)
	}

	return &optionhubproto.GetByNameOut{
		Options: workPlaceList.FromDTO(),
	}, nil
}

func (s *Service) GetWorkPlaceByID(ctx context.Context, in *optionhubproto.GetByIdIn) (*optionhubproto.GetByIdOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetWorkPlaceByID")

	workplace, err := s.dbR.GetWorkPlaceByID(ctx, in.Id)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get workplace by id: %v", err))
		return nil, status.Errorf(codes.NotFound, "failed to get workplace by id: %v", err)
	}

	return &optionhubproto.GetByIdOut{Id: in.Id, Value: workplace}, nil
}

func (s *Service) AddWorkPlace(ctx context.Context, in *optionhubproto.AddIn) (*optionhubproto.AddOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("AddWorkPlace")

	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		logger.Error("failed to find uuid")
		return nil, status.Errorf(codes.Unauthenticated, "failed to find uuid")
	}

	id, err := s.dbR.AddWorkPlace(ctx, in.Value, uuid)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to add new workplace: %v", err))
		return nil, status.Errorf(codes.Aborted, "failed to add new workplace: %v", err)
	}

	return &optionhubproto.AddOut{Id: id, Value: in.Value}, nil
}

func (s *Service) GetStudyPlaceBySearchName(ctx context.Context, in *optionhubproto.GetByNameIn) (*optionhubproto.GetByNameOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetStudyPlaceBySearchName")

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
		logger.Error(fmt.Sprintf("failed to get study place list: %v", err))
		return nil, status.Errorf(codes.NotFound, "failed to get study place list: %v", err)
	}

	return &optionhubproto.GetByNameOut{
		Options: studyPlaceList.FromDTO(),
	}, nil
}

func (s *Service) GetStudyPlaceByID(ctx context.Context, in *optionhubproto.GetByIdIn) (*optionhubproto.GetByIdOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetStudyPlaceByID")

	studyPlace, err := s.dbR.GetStudyPlaceByID(ctx, in.Id)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get study place by id: %v", err))
		return nil, status.Errorf(codes.NotFound, "failed to get study place by id: %v", err)
	}

	return &optionhubproto.GetByIdOut{Id: in.Id, Value: studyPlace}, nil
}

func (s *Service) AddStudyPlace(ctx context.Context, in *optionhubproto.AddIn) (*optionhubproto.AddOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("AddStudyPlace")

	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		logger.Error("failed to find uuid")
		return nil, status.Errorf(codes.Unauthenticated, "failed to find uuid")
	}

	id, err := s.dbR.AddStudyPlace(ctx, in.Value, uuid)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to add new study place: %v", err))
		return nil, status.Errorf(codes.Aborted, "failed to add new study place: %v", err)
	}

	return &optionhubproto.AddOut{Id: id, Value: in.Value}, nil
}

func (s *Service) GetHobbyBySearchName(ctx context.Context, in *optionhubproto.GetByNameIn) (*optionhubproto.GetByNameOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetHobbyBySearchName")

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
		logger.Error(fmt.Sprintf("failed to get hobby list: %v", err))
		return nil, status.Errorf(codes.NotFound, "failed to get hobby list: %v", err)
	}

	return &optionhubproto.GetByNameOut{
		Options: hobbyList.FromDTO(),
	}, nil
}

func (s *Service) GetHobbyByID(ctx context.Context, in *optionhubproto.GetByIdIn) (*optionhubproto.GetByIdOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetHobbyByID")

	hobby, err := s.dbR.GetHobbyByID(ctx, in.Id)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get hobby by id: %v", err))
		return nil, status.Errorf(codes.NotFound, "failed to get hobby by id: %v", err)
	}

	return &optionhubproto.GetByIdOut{Id: in.Id, Value: hobby}, nil
}

func (s *Service) AddHobby(ctx context.Context, in *optionhubproto.AddIn) (*optionhubproto.AddOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("AddHobby")

	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		logger.Error("failed to find uuid")
		return nil, status.Errorf(codes.Unauthenticated, "failed to find uuid")
	}

	id, err := s.dbR.AddHobby(ctx, in.Value, uuid)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to add new hobby: %v", err))
		return nil, status.Errorf(codes.Aborted, "failed to add new hobby: %v", err)
	}

	return &optionhubproto.AddOut{Id: id, Value: in.Value}, nil
}

func (s *Service) GetSkillBySearchName(ctx context.Context, in *optionhubproto.GetByNameIn) (*optionhubproto.GetByNameOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetSkillBySearchName")

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
		logger.Error(fmt.Sprintf("failed to get skill list: %v", err))
		return nil, status.Errorf(codes.NotFound, "failed to get skill list: %v", err)
	}

	return &optionhubproto.GetByNameOut{
		Options: skillList.FromDTO(),
	}, nil
}

func (s *Service) GetSkillByID(ctx context.Context, in *optionhubproto.GetByIdIn) (*optionhubproto.GetByIdOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetSkillByID")

	skill, err := s.dbR.GetSkillByID(ctx, in.Id)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get skill by id: %v", err))
		return nil, status.Errorf(codes.NotFound, "failed to get skill by id: %v", err)
	}

	return &optionhubproto.GetByIdOut{Id: in.Id, Value: skill}, nil
}

func (s *Service) AddSkill(ctx context.Context, in *optionhubproto.AddIn) (*optionhubproto.AddOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("AddSkill")

	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		logger.Error("failed to find uuid")
		return nil, status.Errorf(codes.Unauthenticated, "failed to find uuid")
	}

	id, err := s.dbR.AddSkill(ctx, in.Value, uuid)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to add new skill: %v", err))
		return nil, status.Errorf(codes.Aborted, "failed to add new skill: %v", err)
	}

	return &optionhubproto.AddOut{Id: id, Value: in.Value}, nil
}

func (s *Service) GetCityBySearchName(ctx context.Context, in *optionhubproto.GetByNameIn) (*optionhubproto.GetByNameOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetCityBySearchName")

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
		logger.Error(fmt.Sprintf("failed to get city list: %v", err))
		return nil, status.Errorf(codes.NotFound, "failed to get city list: %v", err)
	}

	return &optionhubproto.GetByNameOut{
		Options: cityList.FromDTO(),
	}, nil
}

func (s *Service) GetCityByID(ctx context.Context, in *optionhubproto.GetByIdIn) (*optionhubproto.GetByIdOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetCityByID")

	city, err := s.dbR.GetCityByID(ctx, in.Id)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get city by id: %v", err))
		return nil, status.Errorf(codes.NotFound, "failed to get city by id: %v", err)
	}

	return &optionhubproto.GetByIdOut{Id: in.Id, Value: city}, nil
}

func (s *Service) AddCity(ctx context.Context, in *optionhubproto.AddIn) (*optionhubproto.AddOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("AddCity")

	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		logger.Error("failed to find uuid")
		return nil, status.Errorf(codes.Unauthenticated, "failed to find uuid")
	}

	id, err := s.dbR.AddCity(ctx, in.Value, uuid)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to add new city: %v", err))
		return nil, status.Errorf(codes.Aborted, "failed to add new city: %v", err)
	}

	return &optionhubproto.AddOut{Id: id, Value: in.Value}, nil
}

func (s *Service) GetSocietyDirectionBySearchName(ctx context.Context, in *optionhubproto.GetByNameIn) (*optionhubproto.GetByNameOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetSocietyDirectionBySearchName")

	var (
		societyDirectionList model.CategoryItemList
		err                  error
	)

	if len(in.Name) < 2 {
		societyDirectionList, err = s.dbR.GetSocietyDirectionPreview(ctx)
	} else {
		societyDirectionList, err = s.dbR.GetSocietyDirectionBySearchName(ctx, in.Name)
	}

	if err != nil {
		logger.Error(fmt.Sprintf("failed to get society direction list: %v", err))
		return nil, status.Errorf(codes.NotFound, "failed to get society direction list: %v", err)
	}

	return &optionhubproto.GetByNameOut{
		Options: societyDirectionList.FromDTO(),
	}, nil
}

func (s *Service) GetSocietyDirectionByID(ctx context.Context, in *optionhubproto.GetByIdIn) (*optionhubproto.GetByIdOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetSocietyDirectionByID")

	societyDirection, err := s.dbR.GetSocietyDirectionByID(ctx, in.Id)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get society direction by id: %v", err))
		return nil, status.Errorf(codes.NotFound, "failed to get society direction by id: %v", err)
	}

	return &optionhubproto.GetByIdOut{Id: in.Id, Value: societyDirection}, nil
}

func (s *Service) AddSocietyDirection(ctx context.Context, in *optionhubproto.AddIn) (*optionhubproto.AddOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("AddSocietyDirection")

	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		logger.Error("failed to find uuid")
		return nil, status.Errorf(codes.Unauthenticated, "failed to find uuid")
	}

	id, err := s.dbR.AddSocietyDirection(ctx, in.Value, uuid)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to add new society direction: %v", err))
		return nil, status.Errorf(codes.Aborted, "failed to add new society direction: %v", err)
	}

	return &optionhubproto.AddOut{Id: id, Value: in.Value}, nil
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
