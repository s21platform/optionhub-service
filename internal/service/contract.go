//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package service

import (
	"context"

	"github.com/s21platform/optionhub-service/internal/model"
)

type DBRepo interface {
	GetAttributesByIds(ctx context.Context, attributesIds []int64) ([]model.Attribute, error)
	GetEntityAttributesByIds(ctx context.Context, attributesEntityIds []int64) ([]model.EntityAttribute, error)
}
