//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package service

import (
	"context"
	"optionhub-service/internal/model"
)

type DBRepo interface {
	AddOS(ctx context.Context, name, uuid string) (int64, error)
	GetOsByID(ctx context.Context, id int64) (string, error)
	GetOsBySearchName(ctx context.Context, name string) (model.CategoryItemList, error)
	GetOsPreview(ctx context.Context) (model.CategoryItemList, error)
	GetAllOs() (model.CategoryItemList, error)

	GetWorkPlaceBySearchName(ctx context.Context, name string) (model.CategoryItemList, error)
	GetWorkPlacePreview(ctx context.Context) (model.CategoryItemList, error)
	GetWorkPlaceByID(ctx context.Context, id int64) (string, error)
	AddWorkPlace(ctx context.Context, name, uuid string) (int64, error)

	GetStudyPlaceBySearchName(ctx context.Context, name string) (model.CategoryItemList, error)
	GetStudyPlacePreview(ctx context.Context) (model.CategoryItemList, error)
	GetStudyPlaceByID(ctx context.Context, id int64) (string, error)
	AddStudyPlace(ctx context.Context, name, uuid string) (int64, error)

	GetHobbyBySearchName(ctx context.Context, name string) (model.CategoryItemList, error)
	GetHobbyPreview(ctx context.Context) (model.CategoryItemList, error)
	GetHobbyByID(ctx context.Context, id int64) (string, error)
	AddHobby(ctx context.Context, name, uuid string) (int64, error)

	GetSkillBySearchName(ctx context.Context, name string) (model.CategoryItemList, error)
	GetSkillPreview(ctx context.Context) (model.CategoryItemList, error)
	GetSkillByID(ctx context.Context, id int64) (string, error)
	AddSkill(ctx context.Context, name, uuid string) (int64, error)

	GetCityBySearchName(ctx context.Context, name string) (model.CategoryItemList, error)
	GetCityPreview(ctx context.Context) (model.CategoryItemList, error)
	GetCityByID(ctx context.Context, id int64) (string, error)
	AddCity(ctx context.Context, name, uuid string) (int64, error)
}
