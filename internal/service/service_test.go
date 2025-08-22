package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/s21platform/optionhub-service/internal/model"
	"github.com/s21platform/optionhub-service/pkg/optionhub"
)

func TestService_GetAttributesMetadata(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		ctx := context.Background()

		mockEntityAttributesIds := []int64{1, 2}
		mockEntityAttributesMeta := []model.EntityAttribute{
			{
				EntityAttributeID: 1,
				EntityType:        "USER",
				AttributeID:       123,
			},
			{
				EntityAttributeID: 2,
				EntityType:        "USER",
				AttributeID:       321,
			},
		}
		mockAttributesIds := []int64{123, 321}
		mockAttributesMeta := []model.Attribute{
			{
				ID:   123,
				Type: "STRING",
			},
			{
				ID:   321,
				Type: "DATE",
			},
		}

		mockDbRepo := NewMockDBRepo(ctrl)
		mockDbRepo.EXPECT().GetEntityAttributesByIds(gomock.Any(), mockEntityAttributesIds).Return(mockEntityAttributesMeta, nil)
		mockDbRepo.EXPECT().GetAttributesByIds(gomock.Any(), mockAttributesIds).Return(mockAttributesMeta, nil)

		s := NewService(mockDbRepo)
		out, err := s.GetAttributesMetadata(ctx, &optionhub.GetAttributesMetadataIn{
			EntityAttributeIds: mockEntityAttributesIds,
		})

		assert.NoError(t, err)
		assert.Equal(t, len(mockEntityAttributesIds), len(out.AttributesMetadata))
	})

	t.Run("fail get attributes metadata", func(t *testing.T) {})
}
