//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package service

import (
	"context"

	"github.com/s21platform/optionhub-service/internal/model"
)

type DBRepo interface {
	GetOptionRequests(ctx context.Context) (model.OptionRequestList, error)
	GetAttributeValueById(ctx context.Context, ids []int64) ([]model.Attribute, error)
	GetValuesByAttributeId(ctx context.Context, attributeId int64) (model.AttributeValueList, error)
	AddAttributeValue(ctx context.Context, in model.AttributeValue) error
}

type SetAttributeProducer interface {
	ProduceMessage(ctx context.Context, message any, key any) error
}
