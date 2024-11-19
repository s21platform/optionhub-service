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
		return 0, fmt.Errorf("cannot execute query, error: %v", err)
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

		return "", fmt.Errorf("cannot execute query, error: %v", err)
	}

	return os, nil
}

// GetOsBySearchName Возвращать то что начинается с name или все совпадения?
func (r *Repository) GetOsBySearchName(ctx context.Context, name string) ([]model.OS, error) {
	var res []model.OS

	searchString := "%" + name + "%"

	query := `SELECT id, name FROM os WHERE name LIKE $1 LIMIT 10`

	err := r.connection.SelectContext(ctx, &res, query, searchString)
	if err != nil {
		return nil, fmt.Errorf("cannot execute query, error: %v", err)
	}

	return res, nil
}
