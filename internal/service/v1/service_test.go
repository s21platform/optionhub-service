package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	logger_lib "github.com/s21platform/logger-lib"
	optionhubproto_v1 "github.com/s21platform/optionhub-proto/optionhub/v1"

	"github.com/s21platform/optionhub-service/internal/config"
	"github.com/s21platform/optionhub-service/internal/model"
	"github.com/s21platform/optionhub-service/utils"
)

func TestService_GetAttributeValues(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	mockRepo := NewMockDBRepo(ctrl)
	kafkaProducer := NewMockSetAttributeProducer(ctrl)

	t.Run("get_attribute_values_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetAttributeValues")

		var attributeId int64 = 5
		expectedDbRes := model.AttributeValueList{
			{
				Id:       1,
				Value:    "Россия",
				ParentId: nil,
			},
			{
				Id:       2,
				Value:    "Москва",
				ParentId: utils.TransformToPtr(int64(1)),
			},
			{
				Id:       3,
				Value:    "Курьяново",
				ParentId: utils.TransformToPtr(int64(2)),
			},
		}

		option3 := &optionhubproto_v1.Option{
			OptionId:    3,
			OptionValue: "Курьяново",
			Children:    []*optionhubproto_v1.Option{},
		}

		option2 := &optionhubproto_v1.Option{
			OptionId:    2,
			OptionValue: "Москва",
			Children:    []*optionhubproto_v1.Option{option3},
		}

		option1 := &optionhubproto_v1.Option{
			OptionId:    1,
			OptionValue: "Россия",
			Children:    []*optionhubproto_v1.Option{option2},
		}

		mockRepo.EXPECT().GetValuesByAttributeId(gomock.Any(), attributeId).Return(expectedDbRes, nil)

		s := NewService(mockRepo, kafkaProducer)
		result, err := s.GetAttributeValues(ctx, &optionhubproto_v1.GetAttributeValuesIn{AttributeId: attributeId})

		assert.NoError(t, err)
		assert.True(t, reflect.DeepEqual(option1, result.OptionList[0]))
	})

	t.Run("get_attribute_values_error", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetAttributeValues")

		expErr := errors.New("cannot get data from db")
		mockLogger.EXPECT().Error(fmt.Sprintf("failed to get attribute values: %v", expErr))

		var attributeId int64 = 5

		mockRepo.EXPECT().GetValuesByAttributeId(gomock.Any(), attributeId).Return(nil, expErr)

		s := NewService(mockRepo, kafkaProducer)
		_, err := s.GetAttributeValues(ctx, &optionhubproto_v1.GetAttributeValuesIn{AttributeId: attributeId})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to get attribute values")
	})
}

func TestService_GetOptionRequests(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	mockRepo := NewMockDBRepo(ctrl)
	kafkaProducer := NewMockSetAttributeProducer(ctrl)

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

		s := NewService(mockRepo, kafkaProducer)
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

		s := NewService(mockRepo, kafkaProducer)
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

		s := NewService(mockRepo, kafkaProducer)
		_, err := s.GetOptionRequests(ctx, &emptypb.Empty{})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to get attribute value by id")
	})
}

func TestService_SetAttribute(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	mockRepo := NewMockDBRepo(ctrl)
	mockProducer := NewMockSetAttributeProducer(ctrl)

	t.Run("set_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("SetAttributeTopic")
		mockRepo.EXPECT().AddAttributeValue(ctx, gomock.Any()).Return(nil)
		mockProducer.EXPECT().ProduceMessage(gomock.Any()).Return(nil)

		s := NewService(mockRepo, mockProducer)
		_, err := s.AddAttributeValue(ctx, &optionhubproto_v1.AddAttributeValueIn{AttributeId: 1, Value: "Linux"})

		assert.NoError(t, err)
	})

	t.Run("set_error", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("SetAttributeTopic")
		mockLogger.EXPECT().Error("failed to add new attribute: test error")

		mockRepo.EXPECT().AddAttributeValue(ctx, gomock.Any()).Return(errors.New("test error"))

		s := NewService(mockRepo, mockProducer)
		_, err := s.AddAttributeValue(ctx, &optionhubproto_v1.AddAttributeValueIn{AttributeId: 1, Value: "Linux"})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Aborted, st.Code())
		assert.Contains(t, st.Message(), "failed to add new attribute")
	})
}
