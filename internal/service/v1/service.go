package service

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	logger_lib "github.com/s21platform/logger-lib"
	new_attribute "github.com/s21platform/optionhub-proto/optionhub/set_new_attribute"
	optionhub "github.com/s21platform/optionhub-proto/optionhub/v1"

	"github.com/s21platform/optionhub-service/internal/config"
	"github.com/s21platform/optionhub-service/internal/model"
)

type Service struct {
	optionhub.UnimplementedOptionhubServiceV1Server
	dbR      DBRepo
	setAttrP SetAttributeProducer
}

func NewService(repo DBRepo, setAttributeProducer SetAttributeProducer) *Service {
	return &Service{dbR: repo, setAttrP: setAttributeProducer}
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

func (s *Service) SetAttribute(ctx context.Context, in *optionhub.SetAttributeByIdIn) error {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("SetAttributeTopic")

	var attributeObj model.AttributeValue

	attributeObj, err := attributeObj.ToDTO(in)

	if err != nil {
		return fmt.Errorf("failed to convert grpc message to dto: %v", err)
	}

	err = s.dbR.SetAttribute(ctx, attributeObj)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to add new attribute: %v", err))
		return status.Errorf(codes.Aborted, "failed to add new attribute: %v", err)
	}

	message := &new_attribute.SetNewAttribute{AttributeId: in.AttributeId}

	err = s.setAttrP.ProduceMessage(message)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to produce kafka message: %v", err))
		return status.Errorf(codes.Aborted, "failed to produce kafka message: %v", err)
	}

	return nil
}
