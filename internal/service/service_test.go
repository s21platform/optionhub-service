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
