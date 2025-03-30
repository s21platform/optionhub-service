//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package service

import (
	"context"

	optionhubproto_v1 "github.com/s21platform/optionhub-proto/optionhub/v1"

	"github.com/s21platform/optionhub-service/internal/model"
)

type DBRepo interface {
	GetOptionRequests(ctx context.Context) (model.OptionRequestList, error)
	GetAttributeValueById(ctx context.Context, ids []int64) ([]model.Attribute, error)

	SetAttribute(ctx context.Context, in *optionhubproto_v1.SetAttributeByIdIn) error
}
