package service_test

import (
	"context"
	"errors"
	"optionhub-service/internal/config"
	"optionhub-service/internal/model/os"
	"optionhub-service/internal/service"
	"testing"

	"github.com/golang/mock/gomock"
	optionhubproto "github.com/s21platform/optionhub-proto/optionhub-proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServer_AddOS(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	uuid := "test-uuid"
	ctx = context.WithValue(ctx, config.KeyUUID, uuid)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := service.NewMockDbRepo(ctrl)

	t.Run("add_ok", func(t *testing.T) {
		osName := "ubuntu"

		var expectedID int64 = 1

		mockRepo.EXPECT().AddOS(gomock.Any(), osName, uuid).Return(expectedID, nil)

		s := service.NewService(mockRepo)
		id, err := s.AddOs(ctx, &optionhubproto.AddIn{Value: osName})
		assert.NoError(t, err)
		assert.Equal(t, id, &optionhubproto.AddOut{Id: expectedID, Value: osName})
	})

	t.Run("add_err", func(t *testing.T) {
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

func TestServer_GetOsById(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := service.NewMockDbRepo(ctrl)

	t.Run("get_by_id_ok", func(t *testing.T) {
		expectedOsName := "ubuntu"

		var id int64 = 3

		mockRepo.EXPECT().GetOsById(gomock.Any(), id).Return(expectedOsName, nil)

		s := service.NewService(mockRepo)
		osName, err := s.GetOsByID(ctx, &optionhubproto.GetByIdIn{Id: id})
		assert.NoError(t, err)
		assert.Equal(t, osName, &optionhubproto.GetByIdOut{Id: id, Value: expectedOsName})
	})

	t.Run("get_by_id_err", func(t *testing.T) {
		var id int64 = 4

		expectedErr := errors.New("get err")

		mockRepo.EXPECT().GetOsById(gomock.Any(), id).Return("", expectedErr)

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

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := service.NewMockDbRepo(ctrl)

	t.Run("get_by_name_ok", func(t *testing.T) {
		expectedNames := []os.Info{
			{ID: 1, Name: "ubuntu"},
			{ID: 2, Name: "ubuntuu"},
			{ID: 5, Name: "ubububu"},
		}
		expectedRes := &optionhubproto.GetByNameOut{
			Values: []*optionhubproto.Record{
				{Id: 1, Value: "ubuntu"},
				{Id: 2, Value: "ubuntuu"},
				{Id: 5, Value: "ubububu"},
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
		search := "w"

		s := service.NewService(mockRepo)
		osNames, err := s.GetOsBySearchName(ctx, &optionhubproto.GetByNameIn{Name: search})
		assert.NoError(t, err)
		assert.Nil(t, osNames)
	})

	t.Run("get_by_name_err", func(t *testing.T) {
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
