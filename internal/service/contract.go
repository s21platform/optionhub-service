//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package service

import (
	"context"
	optionhub_proto "github.com/s21platform/optionhub-proto/optionhub-proto"
)

type DbRepo interface {
	AddOS(ctx context.Context, name string) (int64, error)
	GetOsById(ctx context.Context, id int64) (string, error)
	GetOsBSearchName(ctx context.Context, name string) (*optionhub_proto.GetByNameOut, error)
}
