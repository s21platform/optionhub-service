//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package service

import (
	"context"

	optionhubproto "github.com/s21platform/optionhub-proto/optionhub-proto"
)

type DBRepo interface {
	AddOS(ctx context.Context, name, uuid string) (int64, error)
	GetOsByID(ctx context.Context, id int64) (string, error)
	GetOsBySearchName(ctx context.Context, name string) ([]*optionhubproto.Record, error)
	GetAllOs() ([]*optionhubproto.Record, error)
}
