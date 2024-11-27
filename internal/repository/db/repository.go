package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"optionhub-service/internal/model"
	"time"

	"optionhub-service/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL
)

type Repository struct {
	connection *sqlx.DB
}

func New(cfg *config.Config) (*Repository, error) {
	var err error

	var repo *Repository

	for i := 0; i < 5; i++ {
		repo, err = connect(cfg)
		if err == nil {
			break
		}

		log.Println("failed to connect to database: ", err)
		time.Sleep(500 * time.Millisecond)
	}

	return repo, err
}

func connect(cfg *config.Config) (*Repository, error) {
	conStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
	)

	db, err := sqlx.Connect("postgres", conStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	return &Repository{connection: db}, err
}

func (r *Repository) Close() {
	_ = r.connection.Close()
}

func (r *Repository) AddOS(ctx context.Context, name, uuid string) (int64, error) {
	query := `INSERT INTO os(name, is_moderate, user_uuid) VALUES ($1, $2, $3) RETURNING id`

	var id int64

	err := r.connection.QueryRowxContext(ctx, query, name, true, uuid).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to add os into db: %v", err)
	}

	return id, nil
}

func (r *Repository) GetOsByID(ctx context.Context, id int64) (string, error) {
	var os string

	query := `SELECT name FROM os WHERE id = $1`

	err := r.connection.QueryRowxContext(ctx, query, id).Scan(&os)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", fmt.Errorf("failed to get os by id from db: %v", err)
	}

	return os, nil
}

func (r *Repository) GetOsBySearchName(ctx context.Context, name string) (model.CategoryItemList, error) {
	var res model.CategoryItemList

	searchString := "%" + name + "%"

	query := `SELECT id, name FROM os WHERE name LIKE $1 LIMIT 10`

	err := r.connection.SelectContext(ctx, &res, query, searchString)
	if err != nil {
		return nil, fmt.Errorf("failed to get os by search name from db: %v", err)
	}

	return res, nil
}

func (r *Repository) GetOsPreview(ctx context.Context) (model.CategoryItemList, error) {
	var res model.CategoryItemList

	query := `SELECT id, name FROM os LIMIT 10`

	err := r.connection.SelectContext(ctx, &res, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get os preview from db: %v", err)
	}

	return res, nil
}

func (r *Repository) GetAllOs() (model.CategoryItemList, error) {
	var OSList model.CategoryItemList

	query := `SELECT id, name FROM os`

	err := r.connection.Select(&OSList, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all CategoryItem from db: %w", err)
	}

	return OSList, nil
}

func (r *Repository) GetWorkPlaceBySearchName(ctx context.Context, name string) (model.CategoryItemList, error) {
	var res model.CategoryItemList

	searchString := "%" + name + "%"

	query := `SELECT id, name FROM workplace WHERE name LIKE $1 LIMIT 10`

	err := r.connection.SelectContext(ctx, &res, query, searchString)
	if err != nil {
		return nil, fmt.Errorf("failed to get workplace by search name from db: %v", err)
	}

	return res, nil
}

func (r *Repository) GetWorkPlacePreview(ctx context.Context) (model.CategoryItemList, error) {
	var res model.CategoryItemList

	query := `SELECT id, name FROM workplace LIMIT 10`

	err := r.connection.SelectContext(ctx, &res, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get workplace preview from db: %v", err)
	}

	return res, nil
}

func (r *Repository) GetWorkPlaceByID(ctx context.Context, id int64) (string, error) {
	var workplace string

	query := `SELECT name FROM workplace WHERE id = $1`

	err := r.connection.QueryRowxContext(ctx, query, id).Scan(&workplace)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", fmt.Errorf("failed to get workplace by id from db: %v", err)
	}

	return workplace, nil
}

func (r *Repository) AddWorkPlace(ctx context.Context, name, uuid string) (int64, error) {
	query := `INSERT INTO workplace(name, is_moderate, user_uuid) VALUES ($1, $2, $3) RETURNING id`

	var id int64

	err := r.connection.QueryRowxContext(ctx, query, name, true, uuid).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to add workplace into db: %v", err)
	}

	return id, nil
}

func (r *Repository) GetStudyPlaceBySearchName(ctx context.Context, name string) (model.CategoryItemList, error) {
	var res model.CategoryItemList

	searchString := "%" + name + "%"

	query := `SELECT id, name FROM study_place WHERE name LIKE $1 LIMIT 10`

	err := r.connection.SelectContext(ctx, &res, query, searchString)
	if err != nil {
		return nil, fmt.Errorf("failed to get study place by search name from db: %v", err)
	}

	return res, nil
}

func (r *Repository) GetStudyPlacePreview(ctx context.Context) (model.CategoryItemList, error) {
	var res model.CategoryItemList

	query := `SELECT id, name FROM study_place LIMIT 10`

	err := r.connection.SelectContext(ctx, &res, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get study place preview from db: %v", err)
	}

	return res, nil
}

func (r *Repository) GetStudyPlaceByID(ctx context.Context, id int64) (string, error) {
	var studyPlace string

	query := `SELECT name FROM study_place WHERE id = $1`

	err := r.connection.QueryRowxContext(ctx, query, id).Scan(&studyPlace)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", fmt.Errorf("failed to get study place by id from db: %v", err)
	}

	return studyPlace, nil
}

func (r *Repository) AddStudyPlace(ctx context.Context, name, uuid string) (int64, error) {
	query := `INSERT INTO study_place(name, is_moderate, user_uuid) VALUES ($1, $2, $3) RETURNING id`

	var id int64

	err := r.connection.QueryRowxContext(ctx, query, name, true, uuid).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to add study place into db: %v", err)
	}

	return id, nil
}

func (r *Repository) GetHobbyBySearchName(ctx context.Context, name string) (model.CategoryItemList, error) {
	var res model.CategoryItemList

	searchString := "%" + name + "%"

	query := `SELECT id, name FROM hobby WHERE name LIKE $1 LIMIT 10`

	err := r.connection.SelectContext(ctx, &res, query, searchString)
	if err != nil {
		return nil, fmt.Errorf("failed to get hobby by search name from db: %v", err)
	}

	return res, nil
}

func (r *Repository) GetHobbyPreview(ctx context.Context) (model.CategoryItemList, error) {
	var res model.CategoryItemList

	query := `SELECT id, name FROM hobby LIMIT 10`

	err := r.connection.SelectContext(ctx, &res, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get hobby preview from db: %v", err)
	}

	return res, nil
}

func (r *Repository) GetHobbyByID(ctx context.Context, id int64) (string, error) {
	var hobby string

	query := `SELECT name FROM hobby WHERE id = $1`

	err := r.connection.QueryRowxContext(ctx, query, id).Scan(&hobby)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", fmt.Errorf("failed to get hobby by id from db: %v", err)
	}

	return hobby, nil
}

func (r *Repository) AddHobby(ctx context.Context, name, uuid string) (int64, error) {
	query := `INSERT INTO hobby(name, is_moderate, user_uuid) VALUES ($1, $2, $3) RETURNING id`

	var id int64

	err := r.connection.QueryRowxContext(ctx, query, name, true, uuid).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to add hobby into db: %v", err)
	}

	return id, nil
}

func (r *Repository) GetSkillBySearchName(ctx context.Context, name string) (model.CategoryItemList, error) {
	var res model.CategoryItemList

	searchString := "%" + name + "%"

	query := `SELECT id, name FROM skill WHERE name LIKE $1 LIMIT 10`

	err := r.connection.SelectContext(ctx, &res, query, searchString)
	if err != nil {
		return nil, fmt.Errorf("failed to get skill by search name from db: %v", err)
	}

	return res, nil
}

func (r *Repository) GetSkillPreview(ctx context.Context) (model.CategoryItemList, error) {
	var res model.CategoryItemList

	query := `SELECT id, name FROM skill LIMIT 10`

	err := r.connection.SelectContext(ctx, &res, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get skill preview from db: %v", err)
	}

	return res, nil
}

func (r *Repository) GetSkillByID(ctx context.Context, id int64) (string, error) {
	var skill string

	query := `SELECT name FROM skill WHERE id = $1`

	err := r.connection.QueryRowxContext(ctx, query, id).Scan(&skill)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", fmt.Errorf("failed to get skill by id from db: %v", err)
	}

	return skill, nil
}

func (r *Repository) AddSkill(ctx context.Context, name, uuid string) (int64, error) {
	query := `INSERT INTO skill(name, is_moderate, user_uuid) VALUES ($1, $2, $3) RETURNING id`

	var id int64

	err := r.connection.QueryRowxContext(ctx, query, name, true, uuid).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to add skill into db: %v", err)
	}

	return id, nil
}

func (r *Repository) GetCityBySearchName(ctx context.Context, name string) (model.CategoryItemList, error) {
	var res model.CategoryItemList

	searchString := "%" + name + "%"

	query := `SELECT id, name FROM city WHERE name LIKE $1 LIMIT 10`

	err := r.connection.SelectContext(ctx, &res, query, searchString)
	if err != nil {
		return nil, fmt.Errorf("failed to get city by search name from db: %v", err)
	}

	return res, nil
}

func (r *Repository) GetCityPreview(ctx context.Context) (model.CategoryItemList, error) {
	var res model.CategoryItemList

	query := `SELECT id, name FROM city LIMIT 10`

	err := r.connection.SelectContext(ctx, &res, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get city preview from db: %v", err)
	}

	return res, nil
}

func (r *Repository) GetCityByID(ctx context.Context, id int64) (string, error) {
	var city string

	query := `SELECT name FROM city WHERE id = $1`

	err := r.connection.QueryRowxContext(ctx, query, id).Scan(&city)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", fmt.Errorf("failed to get city by id from db: %v", err)
	}

	return city, nil
}

func (r *Repository) AddCity(ctx context.Context, name, uuid string) (int64, error) {
	query := `INSERT INTO city(name, is_moderate, user_uuid) VALUES ($1, $2, $3) RETURNING id`

	var id int64

	err := r.connection.QueryRowxContext(ctx, query, name, true, uuid).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to add city into db: %v", err)
	}

	return id, nil
}
