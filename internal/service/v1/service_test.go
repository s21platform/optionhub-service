package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	logger_lib "github.com/s21platform/logger-lib"

	"github.com/s21platform/optionhub-service/internal/config"
	"github.com/s21platform/optionhub-service/internal/model"
	"github.com/s21platform/optionhub-service/internal/service/v1"
)

func TestService_GetOptionRequests(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	mockRepo := service.NewMockDBRepo(ctrl)

	t.Run("get_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetOptionRequests")

		now := time.Now()
		expectedRequests := model.OptionRequestList{
			{
				ID:             1,
				AttributeID:    100,
				AttributeValue: "Linux",
				Value:          "Ubuntu",
				UserUuid:       "test-uuid",
				CreatedAt:      now,
			},
		}

		mockRepo.EXPECT().GetOptionRequests(gomock.Any()).Return(expectedRequests, nil)
		mockRepo.EXPECT().GetAttributeValueById(gomock.Any(), []int64{100}).Return([]model.Attribute{{ID: 100, Name: "Linux"}}, nil)

		s := service.NewService(mockRepo)
		result, err := s.GetOptionRequests(ctx, &emptypb.Empty{})

		assert.NoError(t, err)
		assert.Equal(t, int64(1), result.OptionRequestItem[0].OptionRequestId)
		assert.Equal(t, int64(100), result.OptionRequestItem[0].AttributeId)
		assert.Equal(t, "Linux", result.OptionRequestItem[0].AttributeValue)
		assert.Equal(t, "Ubuntu", result.OptionRequestItem[0].OptionRequestValue)
		assert.Equal(t, "test-uuid", result.OptionRequestItem[0].UserUuid)
		assert.Equal(t, timestamppb.New(now), result.OptionRequestItem[0].CreatedAt)
	})

	t.Run("get_error", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetOptionRequests")
		mockLogger.EXPECT().Error("failed to get option requests: test error")

		mockRepo.EXPECT().GetOptionRequests(gomock.Any()).Return(nil, errors.New("test error"))

		s := service.NewService(mockRepo)
		_, err := s.GetOptionRequests(ctx, &emptypb.Empty{})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to get option requests")
	})

	t.Run("get_attributes_error", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetOptionRequests")
		mockLogger.EXPECT().Error("failed to get attribute value by id: test error")

		expectedRequests := model.OptionRequestList{
			{
				ID:          1,
				AttributeID: 100,
			},
		}

		mockRepo.EXPECT().GetOptionRequests(gomock.Any()).Return(expectedRequests, nil)
		mockRepo.EXPECT().GetAttributeValueById(gomock.Any(), []int64{100}).Return(nil, errors.New("test error"))

		s := service.NewService(mockRepo)
		_, err := s.GetOptionRequests(ctx, &emptypb.Empty{})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to get attribute value by id")
	})
}
