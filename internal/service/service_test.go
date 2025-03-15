package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	logger_lib "github.com/s21platform/logger-lib"
	optionhubproto "github.com/s21platform/optionhub-proto/optionhub-proto"

	"optionhub-service/internal/config"
	"optionhub-service/internal/model"
	"optionhub-service/internal/service"
)

func TestServer_AddOS(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uuid := "test-uuid"
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	ctx = context.WithValue(ctx, config.KeyUUID, uuid)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("add_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("AddOs")
		osName := "ubuntu"

		var expectedID int64 = 1

		mockRepo.EXPECT().AddOS(gomock.Any(), osName, uuid).Return(expectedID, nil)

		s := service.NewService(mockRepo)
		id, err := s.AddOs(ctx, &optionhubproto.AddIn{Value: osName})
		assert.NoError(t, err)
		assert.Equal(t, id, &optionhubproto.AddOut{Id: expectedID, Value: osName})
	})

	t.Run("add_no_uuid", func(t *testing.T) {
		ctx := context.Background()
		mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		mockLogger.EXPECT().AddFuncName("AddOs")
		mockLogger.EXPECT().Error("failed to find uuid")

		osName := "macOS"

		s := service.NewService(mockRepo)

		_, err := s.AddOs(ctx, &optionhubproto.AddIn{Value: osName})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
		assert.Contains(t, st.Message(), "failed to find uuid")
	})

	t.Run("add_err", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("AddOs")
		mockLogger.EXPECT().Error("failed to add new os, err: insert err")

		osName := "windows"

		var expectedID int64

		expectedErr := errors.New("insert err")

		mockRepo.EXPECT().AddOS(gomock.Any(), osName, uuid).Return(expectedID, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.AddOs(ctx, &optionhubproto.AddIn{Value: osName})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Aborted, st.Code())
		assert.Contains(t, st.Message(), "insert err")
	})
}

func TestServer_GetOsByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_by_id_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetOsByID")
		expectedOsName := "ubuntu"

		var id int64 = 3

		mockRepo.EXPECT().GetOsByID(gomock.Any(), id).Return(expectedOsName, nil)

		s := service.NewService(mockRepo)
		osName, err := s.GetOsByID(ctx, &optionhubproto.GetByIdIn{Id: id})
		assert.NoError(t, err)
		assert.Equal(t, osName, &optionhubproto.GetByIdOut{Id: id, Value: expectedOsName})
	})

	t.Run("get_by_id_err", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetOsByID")
		mockLogger.EXPECT().Error("failed to get os by id, err: get err")

		var id int64 = 4

		expectedErr := errors.New("get err")

		mockRepo.EXPECT().GetOsByID(gomock.Any(), id).Return("", expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetOsByID(ctx, &optionhubproto.GetByIdIn{Id: id})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "get err")
	})
}

func TestServer_GetOsBySearchName(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_by_name_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetOsBySearchName")

		expectedNames := []model.CategoryItem{
			{ID: 1, Label: "ubuntu"},
			{ID: 2, Label: "ubuntuu"},
			{ID: 5, Label: "ubububu"},
		}
		expectedRes := &optionhubproto.GetByNameOut{
			Options: []*optionhubproto.Record{
				{Id: 1, Label: "ubuntu"},
				{Id: 2, Label: "ubuntuu"},
				{Id: 5, Label: "ubububu"},
			},
		}
		search := "ub"

		mockRepo.EXPECT().GetOsBySearchName(gomock.Any(), search).Return(expectedNames, nil)

		s := service.NewService(mockRepo)
		osNames, err := s.GetOsBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})
		assert.NoError(t, err)
		assert.Equal(t, osNames, expectedRes)
	})

	t.Run("get_by_name_too_less_symbol", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetOsBySearchName")

		search := "w"

		expectedPreview := []model.CategoryItem{
			{ID: 1, Label: "windows"},
			{ID: 2, Label: "wsl"},
		}
		expectedRes := &optionhubproto.GetByNameOut{
			Options: []*optionhubproto.Record{
				{Id: 1, Label: "windows"},
				{Id: 2, Label: "wsl"},
			},
		}

		mockRepo.EXPECT().GetOsPreview(gomock.Any()).Return(expectedPreview, nil)

		s := service.NewService(mockRepo)
		osNames, err := s.GetOsBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})
		assert.NoError(t, err)
		assert.Equal(t, osNames, expectedRes)
	})

	t.Run("get_by_name_err", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetOsBySearchName")
		mockLogger.EXPECT().Error("failed to get os by name, err: db err")

		search := "wi"
		expectedErr := errors.New("db err")

		mockRepo.EXPECT().GetOsBySearchName(gomock.Any(), search).Return(nil, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetOsBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "db err")
	})
}

func TestServer_GetAllOs(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_all_os_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetAllOs")

		expectedNames := []model.CategoryItem{
			{ID: 1, Label: "ubuntu"},
			{ID: 2, Label: "Mac OS"},
			{ID: 5, Label: "Windows"},
		}
		expectedRes := &optionhubproto.GetAllOut{
			Values: []*optionhubproto.Record{
				{Id: 1, Label: "ubuntu"},
				{Id: 2, Label: "Mac OS"},
				{Id: 5, Label: "Windows"},
			},
		}

		mockRepo.EXPECT().GetAllOs().Return(expectedNames, nil)

		s := service.NewService(mockRepo)
		osNames, err := s.GetAllOs(ctx, &optionhubproto.EmptyOptionhub{})
		assert.NoError(t, err)
		assert.Equal(t, osNames, expectedRes)
	})

	t.Run("get_all_os_err", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetAllOs")
		mockLogger.EXPECT().Error("failed to get all os list: db err")

		expectedErr := errors.New("db err")

		mockRepo.EXPECT().GetAllOs().Return(nil, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetAllOs(ctx, &optionhubproto.EmptyOptionhub{})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "db err")
	})
}

func TestServer_GetWorkPlaceBySearchName(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_by_name_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetWorkPlaceBySearchName")

		expectedNames := []model.CategoryItem{
			{ID: 1, Label: "avito"},
			{ID: 2, Label: "avitoo"},
			{ID: 5, Label: "avivito"},
		}
		expectedRes := &optionhubproto.GetByNameOut{
			Options: []*optionhubproto.Record{
				{Id: 1, Label: "avito"},
				{Id: 2, Label: "avitoo"},
				{Id: 5, Label: "avivito"},
			},
		}
		search := "av"

		mockRepo.EXPECT().GetWorkPlaceBySearchName(gomock.Any(), search).Return(expectedNames, nil)

		s := service.NewService(mockRepo)
		workPlaceNames, err := s.GetWorkPlaceBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})
		assert.NoError(t, err)
		assert.Equal(t, workPlaceNames, expectedRes)
	})

	t.Run("get_by_name_too_less_symbol", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetWorkPlaceBySearchName")

		search := "w"

		expectedPreview := []model.CategoryItem{
			{ID: 1, Label: "wildberries"},
			{ID: 2, Label: "ozon"},
		}
		expectedRes := &optionhubproto.GetByNameOut{
			Options: []*optionhubproto.Record{
				{Id: 1, Label: "wildberries"},
				{Id: 2, Label: "ozon"},
			},
		}

		mockRepo.EXPECT().GetWorkPlacePreview(gomock.Any()).Return(expectedPreview, nil)

		s := service.NewService(mockRepo)
		workPlaceNames, err := s.GetWorkPlaceBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})
		assert.NoError(t, err)
		assert.Equal(t, workPlaceNames, expectedRes)
	})

	t.Run("get_by_name_err", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetWorkPlaceBySearchName")
		mockLogger.EXPECT().Error("failed to get workplace list: db err")

		search := "wi"
		expectedErr := errors.New("db err")

		mockRepo.EXPECT().GetWorkPlaceBySearchName(gomock.Any(), search).Return(nil, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetWorkPlaceBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "db err")
	})
}

func TestServer_GetWorkPlaceByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_by_id_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetWorkPlaceByID")

		expectedWorkPlaceName := "avito"

		var id int64 = 3

		mockRepo.EXPECT().GetWorkPlaceByID(gomock.Any(), id).Return(expectedWorkPlaceName, nil)

		s := service.NewService(mockRepo)
		workPlaceName, err := s.GetWorkPlaceByID(ctx, &optionhubproto.GetByIdIn{Id: id})
		assert.NoError(t, err)
		assert.Equal(t, workPlaceName, &optionhubproto.GetByIdOut{Id: id, Value: expectedWorkPlaceName})
	})

	t.Run("get_by_id_err", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetWorkPlaceByID")
		mockLogger.EXPECT().Error("failed to get workplace by id: get err")

		var id int64 = 4

		expectedErr := errors.New("get err")

		mockRepo.EXPECT().GetWorkPlaceByID(gomock.Any(), id).Return("", expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetWorkPlaceByID(ctx, &optionhubproto.GetByIdIn{Id: id})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "get err")
	})
}

func TestServer_AddWorkPlace(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := service.NewMockDBRepo(ctrl)
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)

	ctxWithUUID := context.Background()
	ctxWithUUID = context.WithValue(ctxWithUUID, config.KeyUUID, "test-uuid")
	ctxWithUUID = context.WithValue(ctxWithUUID, config.KeyLogger, mockLogger)

	ctxNoUUID := context.Background()
	ctxNoUUID = context.WithValue(ctxNoUUID, config.KeyLogger, mockLogger)

	t.Run("add_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("AddWorkPlace")

		workPlaceName := "avito"
		var expectedID int64 = 1

		mockRepo.EXPECT().AddWorkPlace(gomock.Any(), workPlaceName, "test-uuid").Return(expectedID, nil)

		s := service.NewService(mockRepo)
		id, err := s.AddWorkPlace(ctxWithUUID, &optionhubproto.AddIn{Value: workPlaceName})
		assert.NoError(t, err)
		assert.Equal(t, id, &optionhubproto.AddOut{Id: expectedID, Value: workPlaceName})
	})

	t.Run("add_no_uuid", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("AddWorkPlace")
		mockLogger.EXPECT().Error("failed to find uuid")

		workPlaceName := "wildberries"

		s := service.NewService(mockRepo)
		_, err := s.AddWorkPlace(ctxNoUUID, &optionhubproto.AddIn{Value: workPlaceName})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
		assert.Contains(t, st.Message(), "failed to find uuid")
	})

	t.Run("add_err", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("AddWorkPlace")
		mockLogger.EXPECT().Error("failed to add new workplace: insert err")

		workPlaceName := "ozon"
		var expectedID int64
		expectedErr := errors.New("insert err")

		mockRepo.EXPECT().AddWorkPlace(gomock.Any(), workPlaceName, "test-uuid").Return(expectedID, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.AddWorkPlace(ctxWithUUID, &optionhubproto.AddIn{Value: workPlaceName})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Aborted, st.Code())
		assert.Contains(t, st.Message(), "insert err")
	})
}

func TestServer_GetStudyPlaceBySearchName(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_by_name_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetStudyPlaceBySearchName")

		expectedNames := []model.CategoryItem{
			{ID: 1, Label: "FU"},
			{ID: 2, Label: "HSE"},
			{ID: 5, Label: "MGIMO"},
		}
		expectedRes := &optionhubproto.GetByNameOut{
			Options: []*optionhubproto.Record{
				{Id: 1, Label: "FU"},
				{Id: 2, Label: "HSE"},
				{Id: 5, Label: "MGIMO"},
			},
		}
		search := "HS"

		mockRepo.EXPECT().GetStudyPlaceBySearchName(gomock.Any(), search).Return(expectedNames, nil)

		s := service.NewService(mockRepo)
		studyPlaceNames, err := s.GetStudyPlaceBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})
		assert.NoError(t, err)
		assert.Equal(t, studyPlaceNames, expectedRes)
	})

	t.Run("get_by_name_too_less_symbol", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetStudyPlaceBySearchName")

		search := "m"

		expectedPreview := []model.CategoryItem{
			{ID: 1, Label: "mirea"},
			{ID: 2, Label: "mgimo"},
		}
		expectedRes := &optionhubproto.GetByNameOut{
			Options: []*optionhubproto.Record{
				{Id: 1, Label: "mirea"},
				{Id: 2, Label: "mgimo"},
			},
		}

		mockRepo.EXPECT().GetStudyPlacePreview(gomock.Any()).Return(expectedPreview, nil)

		s := service.NewService(mockRepo)
		studyPlaceNames, err := s.GetStudyPlaceBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})
		assert.NoError(t, err)
		assert.Equal(t, studyPlaceNames, expectedRes)
	})

	t.Run("get_by_name_err", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetStudyPlaceBySearchName")
		mockLogger.EXPECT().Error("failed to get study place list: db err")

		search := "hs"
		expectedErr := errors.New("db err")

		mockRepo.EXPECT().GetStudyPlaceBySearchName(gomock.Any(), search).Return(nil, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetStudyPlaceBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "db err")
	})
}

func TestServer_GetStudyPlaceByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_by_id_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetStudyPlaceByID")

		expectedStudyPlaceName := "MAI"

		var id int64 = 3

		mockRepo.EXPECT().GetStudyPlaceByID(gomock.Any(), id).Return(expectedStudyPlaceName, nil)

		s := service.NewService(mockRepo)
		studyPlaceName, err := s.GetStudyPlaceByID(ctx, &optionhubproto.GetByIdIn{Id: id})
		assert.NoError(t, err)
		assert.Equal(t, studyPlaceName, &optionhubproto.GetByIdOut{Id: id, Value: expectedStudyPlaceName})
	})

	t.Run("get_by_id_err", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetStudyPlaceByID")
		mockLogger.EXPECT().Error("failed to get study place by id: get err")

		var id int64 = 4

		expectedErr := errors.New("get err")

		mockRepo.EXPECT().GetStudyPlaceByID(gomock.Any(), id).Return("", expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetStudyPlaceByID(ctx, &optionhubproto.GetByIdIn{Id: id})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "get err")
	})
}

func TestServer_AddStudyPlace(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := service.NewMockDBRepo(ctrl)
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)

	ctxWithUUID := context.Background()
	ctxWithUUID = context.WithValue(ctxWithUUID, config.KeyUUID, "test-uuid")
	ctxWithUUID = context.WithValue(ctxWithUUID, config.KeyLogger, mockLogger)

	ctxNoUUID := context.Background()
	ctxNoUUID = context.WithValue(ctxNoUUID, config.KeyLogger, mockLogger)

	t.Run("add_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("AddStudyPlace")
		studyPlaceName := "rudn"
		var expectedID int64 = 1

		mockRepo.EXPECT().AddStudyPlace(gomock.Any(), studyPlaceName, "test-uuid").Return(expectedID, nil)

		s := service.NewService(mockRepo)
		id, err := s.AddStudyPlace(ctxWithUUID, &optionhubproto.AddIn{Value: studyPlaceName})
		assert.NoError(t, err)
		assert.Equal(t, id, &optionhubproto.AddOut{Id: expectedID, Value: studyPlaceName})
	})

	t.Run("add_no_uuid", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("AddStudyPlace")
		mockLogger.EXPECT().Error("failed to find uuid")

		studyPlaceName := "rudn"

		s := service.NewService(mockRepo)
		_, err := s.AddStudyPlace(ctxNoUUID, &optionhubproto.AddIn{Value: studyPlaceName})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
		assert.Contains(t, st.Message(), "failed to find uuid")
	})

	t.Run("add_err", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("AddStudyPlace")
		mockLogger.EXPECT().Error("failed to add new study place: insert err")

		studyPlaceName := "rudn"
		var expectedID int64
		expectedErr := errors.New("insert err")

		mockRepo.EXPECT().AddStudyPlace(gomock.Any(), studyPlaceName, "test-uuid").Return(expectedID, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.AddStudyPlace(ctxWithUUID, &optionhubproto.AddIn{Value: studyPlaceName})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Aborted, st.Code())
		assert.Contains(t, st.Message(), "insert err")
	})
}

func TestServer_GetHobbyBySearchName(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_by_name_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetHobbyBySearchName")

		expectedNames := []model.CategoryItem{
			{ID: 1, Label: "driving"},
			{ID: 2, Label: "painting"},
			{ID: 5, Label: "swimming"},
		}
		expectedRes := &optionhubproto.GetByNameOut{
			Options: []*optionhubproto.Record{
				{Id: 1, Label: "driving"},
				{Id: 2, Label: "painting"},
				{Id: 5, Label: "swimming"},
			},
		}
		search := "pa"

		mockRepo.EXPECT().GetHobbyBySearchName(gomock.Any(), search).Return(expectedNames, nil)

		s := service.NewService(mockRepo)
		hobbyNames, err := s.GetHobbyBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})
		assert.NoError(t, err)
		assert.Equal(t, hobbyNames, expectedRes)
	})

	t.Run("get_by_name_too_less_symbol", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetHobbyBySearchName")

		search := "t"

		expectedPreview := []model.CategoryItem{
			{ID: 1, Label: "coding"},
			{ID: 2, Label: "testing"},
		}
		expectedRes := &optionhubproto.GetByNameOut{
			Options: []*optionhubproto.Record{
				{Id: 1, Label: "coding"},
				{Id: 2, Label: "testing"},
			},
		}

		mockRepo.EXPECT().GetHobbyPreview(gomock.Any()).Return(expectedPreview, nil)

		s := service.NewService(mockRepo)
		hobbyNames, err := s.GetHobbyBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})
		assert.NoError(t, err)
		assert.Equal(t, hobbyNames, expectedRes)
	})

	t.Run("get_by_name_err", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetHobbyBySearchName")
		mockLogger.EXPECT().Error("failed to get hobby list: db err")

		search := "dr"
		expectedErr := errors.New("db err")

		mockRepo.EXPECT().GetHobbyBySearchName(gomock.Any(), search).Return(nil, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetHobbyBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "db err")
	})
}

func TestServer_GetHobbyByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_by_id_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetHobbyByID")

		expectedHobbyName := "singing"

		var id int64 = 3

		mockRepo.EXPECT().GetHobbyByID(gomock.Any(), id).Return(expectedHobbyName, nil)

		s := service.NewService(mockRepo)
		hobbyName, err := s.GetHobbyByID(ctx, &optionhubproto.GetByIdIn{Id: id})
		assert.NoError(t, err)
		assert.Equal(t, hobbyName, &optionhubproto.GetByIdOut{Id: id, Value: expectedHobbyName})
	})

	t.Run("get_by_id_err", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetHobbyByID")
		mockLogger.EXPECT().Error("failed to get hobby by id: get err")

		var id int64 = 4

		expectedErr := errors.New("get err")

		mockRepo.EXPECT().GetHobbyByID(gomock.Any(), id).Return("", expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetHobbyByID(ctx, &optionhubproto.GetByIdIn{Id: id})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "get err")
	})
}

func TestServer_AddHobby(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := service.NewMockDBRepo(ctrl)
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)

	ctxWithUUID := context.Background()
	ctxWithUUID = context.WithValue(ctxWithUUID, config.KeyUUID, "test-uuid")
	ctxWithUUID = context.WithValue(ctxWithUUID, config.KeyLogger, mockLogger)

	ctxNoUUID := context.Background()
	ctxNoUUID = context.WithValue(ctxNoUUID, config.KeyLogger, mockLogger)

	t.Run("add_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("AddHobby")

		hobbyName := "boxing"
		var expectedID int64 = 1

		mockRepo.EXPECT().AddHobby(gomock.Any(), hobbyName, "test-uuid").Return(expectedID, nil)

		s := service.NewService(mockRepo)
		id, err := s.AddHobby(ctxWithUUID, &optionhubproto.AddIn{Value: hobbyName})
		assert.NoError(t, err)
		assert.Equal(t, id, &optionhubproto.AddOut{Id: expectedID, Value: hobbyName})
	})

	t.Run("add_no_uuid", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("AddHobby")
		mockLogger.EXPECT().Error("failed to find uuid")

		hobbyName := "boxing"

		s := service.NewService(mockRepo)
		_, err := s.AddHobby(ctxNoUUID, &optionhubproto.AddIn{Value: hobbyName})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
		assert.Contains(t, st.Message(), "failed to find uuid")
	})

	t.Run("add_err", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("AddHobby")
		mockLogger.EXPECT().Error("failed to add new hobby: insert err")

		hobbyName := "boxing"
		var expectedID int64
		expectedErr := errors.New("insert err")

		mockRepo.EXPECT().AddHobby(gomock.Any(), hobbyName, "test-uuid").Return(expectedID, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.AddHobby(ctxWithUUID, &optionhubproto.AddIn{Value: hobbyName})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Aborted, st.Code())
		assert.Contains(t, st.Message(), "insert err")
	})
}

func TestServer_GetSkillBySearchName(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_by_name_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetSkillBySearchName")

		expectedNames := []model.CategoryItem{
			{ID: 1, Label: "kafka"},
			{ID: 2, Label: "s3"},
			{ID: 5, Label: "swift"},
		}
		expectedRes := &optionhubproto.GetByNameOut{
			Options: []*optionhubproto.Record{
				{Id: 1, Label: "kafka"},
				{Id: 2, Label: "s3"},
				{Id: 5, Label: "swift"},
			},
		}
		search := "sw"

		mockRepo.EXPECT().GetSkillBySearchName(gomock.Any(), search).Return(expectedNames, nil)

		s := service.NewService(mockRepo)
		skillNames, err := s.GetSkillBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})
		assert.NoError(t, err)
		assert.Equal(t, skillNames, expectedRes)
	})

	t.Run("get_by_name_too_less_symbol", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetSkillBySearchName")

		search := "q"

		expectedPreview := []model.CategoryItem{
			{ID: 1, Label: "sql"},
			{ID: 2, Label: "python"},
		}
		expectedRes := &optionhubproto.GetByNameOut{
			Options: []*optionhubproto.Record{
				{Id: 1, Label: "sql"},
				{Id: 2, Label: "python"},
			},
		}

		mockRepo.EXPECT().GetSkillPreview(gomock.Any()).Return(expectedPreview, nil)

		s := service.NewService(mockRepo)
		skillNames, err := s.GetSkillBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})
		assert.NoError(t, err)
		assert.Equal(t, skillNames, expectedRes)
	})

	t.Run("get_by_name_err", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetSkillBySearchName")
		mockLogger.EXPECT().Error("failed to get skill list: db err")

		search := "go"
		expectedErr := errors.New("db err")

		mockRepo.EXPECT().GetSkillBySearchName(gomock.Any(), search).Return(nil, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetSkillBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "db err")
	})
}

func TestServer_GetSkillByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_by_id_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetSkillByID")

		expectedSkillName := "kotlin"

		var id int64 = 3

		mockRepo.EXPECT().GetSkillByID(gomock.Any(), id).Return(expectedSkillName, nil)

		s := service.NewService(mockRepo)
		skillName, err := s.GetSkillByID(ctx, &optionhubproto.GetByIdIn{Id: id})
		assert.NoError(t, err)
		assert.Equal(t, skillName, &optionhubproto.GetByIdOut{Id: id, Value: expectedSkillName})
	})

	t.Run("get_by_id_err", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetSkillByID")
		mockLogger.EXPECT().Error("failed to get skill by id: get err")

		var id int64 = 4

		expectedErr := errors.New("get err")

		mockRepo.EXPECT().GetSkillByID(gomock.Any(), id).Return("", expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetSkillByID(ctx, &optionhubproto.GetByIdIn{Id: id})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "get err")
	})
}

func TestServer_AddSkill(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := service.NewMockDBRepo(ctrl)
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)

	ctxWithUUID := context.Background()
	ctxWithUUID = context.WithValue(ctxWithUUID, config.KeyUUID, "test-uuid")
	ctxWithUUID = context.WithValue(ctxWithUUID, config.KeyLogger, mockLogger)

	ctxNoUUID := context.Background()
	ctxNoUUID = context.WithValue(ctxNoUUID, config.KeyLogger, mockLogger)

	t.Run("add_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("AddSkill")

		skillName := "R"
		var expectedID int64 = 1

		mockRepo.EXPECT().AddSkill(gomock.Any(), skillName, "test-uuid").Return(expectedID, nil)

		s := service.NewService(mockRepo)
		id, err := s.AddSkill(ctxWithUUID, &optionhubproto.AddIn{Value: skillName})
		assert.NoError(t, err)
		assert.Equal(t, id, &optionhubproto.AddOut{Id: expectedID, Value: skillName})
	})

	t.Run("add_no_uuid", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("AddSkill")
		mockLogger.EXPECT().Error("failed to find uuid")

		skillName := "R"

		s := service.NewService(mockRepo)
		_, err := s.AddSkill(ctxNoUUID, &optionhubproto.AddIn{Value: skillName})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
		assert.Contains(t, st.Message(), "failed to find uuid")
	})

	t.Run("add_err", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("AddSkill")
		mockLogger.EXPECT().Error("failed to add new skill: insert err")

		skillName := "R"
		var expectedID int64
		expectedErr := errors.New("insert err")

		mockRepo.EXPECT().AddSkill(gomock.Any(), skillName, "test-uuid").Return(expectedID, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.AddSkill(ctxWithUUID, &optionhubproto.AddIn{Value: skillName})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Aborted, st.Code())
		assert.Contains(t, st.Message(), "insert err")
	})
}

func TestServer_GetCityBySearchName(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_by_name_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetCityBySearchName")

		expectedNames := []model.CategoryItem{
			{ID: 1, Label: "Moscow"},
			{ID: 2, Label: "New York"},
			{ID: 5, Label: "St. Petersburg"},
		}
		expectedRes := &optionhubproto.GetByNameOut{
			Options: []*optionhubproto.Record{
				{Id: 1, Label: "Moscow"},
				{Id: 2, Label: "New York"},
				{Id: 5, Label: "St. Petersburg"},
			},
		}
		search := "Mo"

		mockRepo.EXPECT().GetCityBySearchName(gomock.Any(), search).Return(expectedNames, nil)

		s := service.NewService(mockRepo)
		cityNames, err := s.GetCityBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})
		assert.NoError(t, err)
		assert.Equal(t, cityNames, expectedRes)
	})

	t.Run("get_by_name_too_less_symbol", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetCityBySearchName")

		search := "v"

		expectedPreview := []model.CategoryItem{
			{ID: 1, Label: "Voronezh"},
			{ID: 2, Label: "Vena"},
		}
		expectedRes := &optionhubproto.GetByNameOut{
			Options: []*optionhubproto.Record{
				{Id: 1, Label: "Voronezh"},
				{Id: 2, Label: "Vena"},
			},
		}

		mockRepo.EXPECT().GetCityPreview(gomock.Any()).Return(expectedPreview, nil)

		s := service.NewService(mockRepo)
		cityNames, err := s.GetCityBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})
		assert.NoError(t, err)
		assert.Equal(t, cityNames, expectedRes)
	})

	t.Run("get_by_name_err", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetCityBySearchName")
		mockLogger.EXPECT().Error("failed to get city list: db err")

		search := "ne"
		expectedErr := errors.New("db err")

		mockRepo.EXPECT().GetCityBySearchName(gomock.Any(), search).Return(nil, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetCityBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "db err")
	})
}

func TestServer_GetCityByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_by_id_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetCityByID")

		expectedCityName := "Almata"

		var id int64 = 3

		mockRepo.EXPECT().GetCityByID(gomock.Any(), id).Return(expectedCityName, nil)

		s := service.NewService(mockRepo)
		cityName, err := s.GetCityByID(ctx, &optionhubproto.GetByIdIn{Id: id})
		assert.NoError(t, err)
		assert.Equal(t, cityName, &optionhubproto.GetByIdOut{Id: id, Value: expectedCityName})
	})

	t.Run("get_by_id_err", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetCityByID")
		mockLogger.EXPECT().Error("failed to get city by id: get err")

		var id int64 = 4

		expectedErr := errors.New("get err")

		mockRepo.EXPECT().GetCityByID(gomock.Any(), id).Return("", expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetCityByID(ctx, &optionhubproto.GetByIdIn{Id: id})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "get err")
	})
}

func TestServer_AddCity(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := service.NewMockDBRepo(ctrl)
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)

	ctxWithUUID := context.Background()
	ctxWithUUID = context.WithValue(ctxWithUUID, config.KeyUUID, "test-uuid")
	ctxWithUUID = context.WithValue(ctxWithUUID, config.KeyLogger, mockLogger)

	ctxNoUUID := context.Background()
	ctxNoUUID = context.WithValue(ctxNoUUID, config.KeyLogger, mockLogger)

	t.Run("add_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("AddCity")

		cityName := "Dublin"
		var expectedID int64 = 1

		mockRepo.EXPECT().AddCity(gomock.Any(), cityName, "test-uuid").Return(expectedID, nil)

		s := service.NewService(mockRepo)
		id, err := s.AddCity(ctxWithUUID, &optionhubproto.AddIn{Value: cityName})
		assert.NoError(t, err)
		assert.Equal(t, id, &optionhubproto.AddOut{Id: expectedID, Value: cityName})
	})

	t.Run("add_no_uuid", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("AddCity")
		mockLogger.EXPECT().Error("failed to find uuid")

		cityName := "Dublin"

		s := service.NewService(mockRepo)
		_, err := s.AddCity(ctxNoUUID, &optionhubproto.AddIn{Value: cityName})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
		assert.Contains(t, st.Message(), "failed to find uuid")
	})

	t.Run("add_err", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("AddCity")
		mockLogger.EXPECT().Error("failed to add new city: insert err")

		cityName := "Dublin"
		var expectedID int64
		expectedErr := errors.New("insert err")

		mockRepo.EXPECT().AddCity(gomock.Any(), cityName, "test-uuid").Return(expectedID, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.AddCity(ctxWithUUID, &optionhubproto.AddIn{Value: cityName})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Aborted, st.Code())
		assert.Contains(t, st.Message(), "insert err")
	})
}

func TestServer_GetSocietyDirectionBySearchName(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_by_name_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetSocietyDirectionBySearchName")

		expectedNames := []model.CategoryItem{
			{ID: 1, Label: "Cybersecurity"},
			{ID: 2, Label: "E-Government"},
			{ID: 5, Label: "EdTech"},
		}
		expectedRes := &optionhubproto.GetByNameOut{
			Options: []*optionhubproto.Record{
				{Id: 1, Label: "Cybersecurity"},
				{Id: 2, Label: "E-Government"},
				{Id: 5, Label: "EdTech"},
			},
		}
		search := "Cy"

		mockRepo.EXPECT().GetSocietyDirectionBySearchName(gomock.Any(), search).Return(expectedNames, nil)

		s := service.NewService(mockRepo)
		societyDirectionNames, err := s.GetSocietyDirectionBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})
		assert.NoError(t, err)
		assert.Equal(t, societyDirectionNames, expectedRes)
	})

	t.Run("get_by_name_too_less_symbol", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetSocietyDirectionBySearchName")

		search := "A"

		expectedPreview := []model.CategoryItem{
			{ID: 1, Label: "AI"},
			{ID: 2, Label: "Blockchain"},
		}
		expectedRes := &optionhubproto.GetByNameOut{
			Options: []*optionhubproto.Record{
				{Id: 1, Label: "AI"},
				{Id: 2, Label: "Blockchain"},
			},
		}

		mockRepo.EXPECT().GetSocietyDirectionPreview(gomock.Any()).Return(expectedPreview, nil)

		s := service.NewService(mockRepo)
		societyDirectionNames, err := s.GetSocietyDirectionBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})
		assert.NoError(t, err)
		assert.Equal(t, societyDirectionNames, expectedRes)
	})

	t.Run("get_by_name_err", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetSocietyDirectionBySearchName")
		mockLogger.EXPECT().Error("failed to get society direction list: db err")

		search := "ta"
		expectedErr := errors.New("db err")

		mockRepo.EXPECT().GetSocietyDirectionBySearchName(gomock.Any(), search).Return(nil, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetSocietyDirectionBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "db err")
	})
}

func TestServer_GetSocietyDirectionByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_by_id_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetSocietyDirectionByID")

		expectedSocietyDirectionName := "Sustainability"

		var id int64 = 3

		mockRepo.EXPECT().GetSocietyDirectionByID(gomock.Any(), id).Return(expectedSocietyDirectionName, nil)

		s := service.NewService(mockRepo)
		societyDirectionName, err := s.GetSocietyDirectionByID(ctx, &optionhubproto.GetByIdIn{Id: id})
		assert.NoError(t, err)
		assert.Equal(t, societyDirectionName, &optionhubproto.GetByIdOut{Id: id, Value: expectedSocietyDirectionName})
	})

	t.Run("get_by_id_err", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetSocietyDirectionByID")
		mockLogger.EXPECT().Error("failed to get society direction by id: get err")

		var id int64 = 4

		expectedErr := errors.New("get err")

		mockRepo.EXPECT().GetSocietyDirectionByID(gomock.Any(), id).Return("", expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.GetSocietyDirectionByID(ctx, &optionhubproto.GetByIdIn{Id: id})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "get err")
	})
}

func TestServer_AddSocietyDirection(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := service.NewMockDBRepo(ctrl)
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)

	ctxWithUUID := context.Background()
	ctxWithUUID = context.WithValue(ctxWithUUID, config.KeyUUID, "test-uuid")
	ctxWithUUID = context.WithValue(ctxWithUUID, config.KeyLogger, mockLogger)

	ctxNoUUID := context.Background()
	ctxNoUUID = context.WithValue(ctxNoUUID, config.KeyLogger, mockLogger)

	t.Run("add_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("AddSocietyDirection")

		societyDirectionName := "Privacy"
		var expectedID int64 = 1

		mockRepo.EXPECT().AddSocietyDirection(gomock.Any(), societyDirectionName, "test-uuid").Return(expectedID, nil)

		s := service.NewService(mockRepo)
		id, err := s.AddSocietyDirection(ctxWithUUID, &optionhubproto.AddIn{Value: societyDirectionName})
		assert.NoError(t, err)
		assert.Equal(t, id, &optionhubproto.AddOut{Id: expectedID, Value: societyDirectionName})
	})

	t.Run("add_no_uuid", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("AddSocietyDirection")
		mockLogger.EXPECT().Error("failed to find uuid")

		societyDirectionName := "Privacy"

		s := service.NewService(mockRepo)
		_, err := s.AddSocietyDirection(ctxNoUUID, &optionhubproto.AddIn{Value: societyDirectionName})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
		assert.Contains(t, st.Message(), "failed to find uuid")
	})

	t.Run("add_err", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("AddSocietyDirection")
		mockLogger.EXPECT().Error("failed to add new society direction: insert err")

		societyDirectionName := "Privacy"
		var expectedID int64
		expectedErr := errors.New("insert err")

		mockRepo.EXPECT().AddSocietyDirection(gomock.Any(), societyDirectionName, "test-uuid").Return(expectedID, expectedErr)

		s := service.NewService(mockRepo)
		_, err := s.AddSocietyDirection(ctxWithUUID, &optionhubproto.AddIn{Value: societyDirectionName})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Aborted, st.Code())
		assert.Contains(t, st.Message(), "insert err")
	})
}
