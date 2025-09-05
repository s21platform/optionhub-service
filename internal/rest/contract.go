package rest

import (
	"context"

	"github.com/s21platform/optionhub-service/internal/model"
)

type DbRepo interface {
	GetOptionsByAttributeId(ctx context.Context, attributeId int64) ([]model.Option, error)
}