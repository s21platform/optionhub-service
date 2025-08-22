package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	logger_lib "github.com/s21platform/logger-lib"

	"github.com/s21platform/optionhub-service/internal/model"
	"github.com/s21platform/optionhub-service/pkg/optionhub"
)

type Service struct {
	optionhub.UnimplementedOptionhubServiceServer
	dbR DBRepo
}

func NewService(repo DBRepo) *Service {
	return &Service{dbR: repo}
}

func (s *Service) GetAttributesMetadata(ctx context.Context, in *optionhub.GetAttributesMetadataIn) (*optionhub.GetAttributesMetadataOut, error) {
	if len(in.EntityAttributeIds) == 0 {
		return nil, status.Error(codes.InvalidArgument, "entity attribute ids is required")
	}

	entityAttributes, err := s.dbR.GetEntityAttributesByIds(ctx, in.EntityAttributeIds)
	if err != nil {
		logger_lib.Error(logger_lib.WithField(ctx, "error", err), "failed to get entity attributes")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "entity attributes not found")
		}
		return nil, status.Error(codes.Internal, "failed to get entity attributes")
	}

	attributeIds := lo.Map(entityAttributes, func(item model.EntityAttribute, index int) int64 {
		return item.AttributeID
	})

	attributes, err := s.dbR.GetAttributesByIds(ctx, attributeIds)
	if err != nil {
		logger_lib.Error(logger_lib.WithField(ctx, "error", err), "failed to get attributes")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "attributes not found")
		}
		return nil, status.Error(codes.Internal, "failed to get attributes")
	}

	attrMetaMap := make(map[int64]model.Attribute, len(attributes))
	for _, meta := range attributes {
		attrMetaMap[meta.ID] = meta
	}

	res := make(map[int64]*optionhub.AttributeMetadata, len(entityAttributes))
	for _, entityAttr := range entityAttributes {
		if attrMeta, ok := attrMetaMap[entityAttr.AttributeID]; ok {
			res[entityAttr.EntityAttributeID] = &optionhub.AttributeMetadata{
				EntityAttributeId: entityAttr.EntityAttributeID,
				AttributeId:       entityAttr.AttributeID,
				Type:              optionhub.AttributeType(optionhub.AttributeType_value[attrMeta.Type]),
				Name:              attrMeta.Name,
				Description:       attrMeta.Description,
				EntityType:        optionhub.EntityType(optionhub.EntityType_value[entityAttr.EntityType]),
				Label:             entityAttr.Label,
				IsRequired:        entityAttr.IsRequired,
				OrderIndex:        int64(entityAttr.OrderIndex),
				VisibilityRules:   entityAttr.VisibilityRules,
			}
		}
	}

	return &optionhub.GetAttributesMetadataOut{
		AttributesMetadata: res,
	}, nil
}
