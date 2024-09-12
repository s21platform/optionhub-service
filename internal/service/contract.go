package service

import "context"

type DbRepo interface {
	AddOS(ctx context.Context, name string) (int64, error)
	GetOsById(ctx context.Context, id int64) (string, error)
}
