package service

import (
	"context"
	"database/sql"
	"errors"
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

	t.Run("fail no attributes in", func(t *testing.T) {
		ctx := context.Background()

		s := NewService(nil)

		_, err := s.GetAttributesMetadata(ctx, &optionhub.GetAttributesMetadataIn{})
		assert.Error(t, err)
	})

	t.Run("fail to get entity attributes no rows", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		ctx := context.Background()

		mockEntityAttributesIds := []int64{1, 2}

		mockDbRepo := NewMockDBRepo(ctrl)
		mockDbRepo.EXPECT().GetEntityAttributesByIds(gomock.Any(), mockEntityAttributesIds).Return(nil, sql.ErrNoRows)

		s := NewService(mockDbRepo)
		_, err := s.GetAttributesMetadata(ctx, &optionhub.GetAttributesMetadataIn{
			EntityAttributeIds: mockEntityAttributesIds,
		})
		assert.Error(t, err)
	})

	t.Run("fail to get entity attributes err", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		ctx := context.Background()

		mockEntityAttributesIds := []int64{1, 2}

		mockDbRepo := NewMockDBRepo(ctrl)
		mockDbRepo.EXPECT().GetEntityAttributesByIds(gomock.Any(), mockEntityAttributesIds).Return(nil, errors.New("some error"))

		s := NewService(mockDbRepo)
		_, err := s.GetAttributesMetadata(ctx, &optionhub.GetAttributesMetadataIn{
			EntityAttributeIds: mockEntityAttributesIds,
		})
		assert.Error(t, err)
	})

	t.Run("fail to get attributes no rows", func(t *testing.T) {
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

		mockDbRepo := NewMockDBRepo(ctrl)
		mockDbRepo.EXPECT().GetEntityAttributesByIds(gomock.Any(), mockEntityAttributesIds).Return(mockEntityAttributesMeta, nil)
		mockDbRepo.EXPECT().GetAttributesByIds(gomock.Any(), mockAttributesIds).Return(nil, sql.ErrNoRows)

		s := NewService(mockDbRepo)
		_, err := s.GetAttributesMetadata(ctx, &optionhub.GetAttributesMetadataIn{
			EntityAttributeIds: mockEntityAttributesIds,
		})
		assert.Error(t, err)
	})

	t.Run("fail to get attributes err", func(t *testing.T) {
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

		mockDbRepo := NewMockDBRepo(ctrl)
		mockDbRepo.EXPECT().GetEntityAttributesByIds(gomock.Any(), mockEntityAttributesIds).Return(mockEntityAttributesMeta, nil)
		mockDbRepo.EXPECT().GetAttributesByIds(gomock.Any(), mockAttributesIds).Return(nil, errors.New("some error"))

		s := NewService(mockDbRepo)
		_, err := s.GetAttributesMetadata(ctx, &optionhub.GetAttributesMetadataIn{
			EntityAttributeIds: mockEntityAttributesIds,
		})
		assert.Error(t, err)
	})
}
