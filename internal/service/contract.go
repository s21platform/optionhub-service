//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package service

import "context"

type DbRepo interface {
	AddOS(ctx context.Context, name string) (int64, error)
	GetOsById(ctx context.Context, id int64) (string, error)
}
