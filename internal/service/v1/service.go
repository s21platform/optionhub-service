package service

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	logger_lib "github.com/s21platform/logger-lib"
	optionhub "github.com/s21platform/optionhub-proto/optionhub/v1"

	"github.com/s21platform/optionhub-service/internal/config"
	"github.com/s21platform/optionhub-service/internal/model"
)

type Service struct {
	optionhub.UnimplementedOptionhubServiceV1Server
	dbR DBRepo
}

func NewService(repo DBRepo) *Service {
	return &Service{dbR: repo}
}

func (s *Service) GetAttributeValues(ctx context.Context, in *optionhub.GetAttributeValuesIn) (*optionhub.GetAttributeValuesOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetAttributeValues")

	values, err := s.dbR.GetValuesByAttributeId(ctx, in.AttributeId)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get attribute values: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to get attribute values: %v", err)
	}

	return &optionhub.GetAttributeValuesOut{OptionList: values.ToDTO()}, nil
}

func (s *Service) GetOptionRequests(ctx context.Context, _ *emptypb.Empty) (*optionhub.GetOptionRequestsOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetOptionRequests")

	requests, err := s.dbR.GetOptionRequests(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get option requests: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to get option requests: %v", err)
	}

	attributes, err := s.dbR.GetAttributeValueById(ctx, lo.Map(requests, func(o model.OptionRequest, _ int) int64 { return o.AttributeID }))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get attribute value by id: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to get attribute value by id: %v", err)
	}

	resp := requests.ToDTO()

	attributeMap := lo.KeyBy(attributes, func(a model.Attribute) int64 { return a.ID })
	lo.ForEach(resp, func(o *optionhub.OptionRequestItem, _ int) {
		if attr, ok := attributeMap[o.AttributeId]; ok {
			o.AttributeValue = attr.Name
		}
	})

	return &optionhub.GetOptionRequestsOut{
		OptionRequestItem: resp,
	}, nil
}
