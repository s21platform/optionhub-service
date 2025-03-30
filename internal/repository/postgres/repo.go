package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL

	"github.com/s21platform/optionhub-service/internal/config"
	"github.com/s21platform/optionhub-service/internal/model"
)

type Repository struct {
	connection *sqlx.DB
}

func New(cfg *config.Config) *Repository {
	conStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database, cfg.Postgres.Host, cfg.Postgres.Port)

	conn, err := sqlx.Connect("postgres", conStr)
	if err != nil {
		log.Fatal("error connect: ", err)
	}

	return &Repository{
		connection: conn,
	}
}

func (r *Repository) Close() {
	_ = r.connection.Close()
}

func (r *Repository) AddOS(ctx context.Context, name, uuid string) (int64, error) {
	query := `INSERT INTO os(name, is_moderate, user_uuid) VALUES ($1, $2, $3) RETURNING id`

	var id int64

	err := r.connection.QueryRowxContext(ctx, query, name, true, uuid).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to add os into postgres: %v", err)
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

		return "", fmt.Errorf("failed to get os by id from postgres: %v", err)
	}

	return os, nil
}

func (r *Repository) GetOsBySearchName(ctx context.Context, name string) (model.CategoryItemList, error) {
	var res model.CategoryItemList

	searchString := "%" + name + "%"

	query := `SELECT id, name FROM os WHERE name ILIKE $1 LIMIT 10`

	err := r.connection.SelectContext(ctx, &res, query, searchString)
	if err != nil {
		return nil, fmt.Errorf("failed to get os by search name from postgres: %v", err)
	}

	return res, nil
}

func (r *Repository) GetOsPreview(ctx context.Context) (model.CategoryItemList, error) {
	var res model.CategoryItemList

	query := `SELECT id, name FROM os LIMIT 10`

	err := r.connection.SelectContext(ctx, &res, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get os preview from postgres: %v", err)
	}

	return res, nil
}

func (r *Repository) GetAllOs() (model.CategoryItemList, error) {
	var OSList model.CategoryItemList

	query := `SELECT id, name FROM os`

	err := r.connection.Select(&OSList, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all CategoryItem from postgres: %w", err)
	}

	return OSList, nil
}

func (r *Repository) GetAttributeValueById(ctx context.Context, ids []int64) ([]model.Attribute, error) {
	var res []model.Attribute

	query, args, err := sq.
		Select("id", "name").
		From("attributes").
		Where(sq.Eq{"id": ids}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	err = r.connection.SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	return res, nil
}

func (r *Repository) GetOptionRequests(ctx context.Context) (model.OptionRequestList, error) {
	var res model.OptionRequestList

	query, args, err := sq.
		Select(
			"id",
			"attribute_id",
			"value",
			"user_uuid",
			"created_at",
		).
		From("option_requests").
		OrderBy("id DESC").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	err = r.connection.SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get option requests: %v", err)
	}

	return res, nil
}
