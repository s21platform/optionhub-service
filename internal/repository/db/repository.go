package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"optionhub-service/internal/model/os"
	"time"

	"optionhub-service/internal/config"

	_ "github.com/lib/pq" // for postgres
)

type Repository struct {
	connection *sql.DB
}

func connect(cfg *config.Config) (*Repository, error) {
	// Connect db
	conStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database, cfg.Postgres.Host, cfg.Postgres.Port)

	db, err := sql.Open("postgres", conStr)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	// Сhecking connection db
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping: %w", err)
	}

	return &Repository{db}, nil
}

func (r *Repository) Close() {
	r.connection.Close()
}

func New(cfg *config.Config) (*Repository, error) {
	var err error

	var repo *Repository

	for i := 0; i < 5; i++ {
		repo, err = connect(cfg)
		if err == nil {
			return repo, nil
		}

		log.Println(err)
		time.Sleep(500 * time.Millisecond)
	}

	return nil, err
}

func (r *Repository) AddOS(ctx context.Context, name, uuid string) (int64, error) {
	query := "INSERT INTO os(name, create_at, is_moderate, user_uuid) VALUES ($1, $2, $3, $4) RETURNING id"
	createTime := time.Now().UTC()

	var id int64

	err := r.connection.QueryRowContext(ctx, query, name, createTime, true, uuid).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("cannot execute query, error: %v", err)
	}

	return id, nil
}

func (r *Repository) GetOsByID(ctx context.Context, id int64) (string, error) {
	var os string

	query := "SELECT name from os where id = $1"

	err := r.connection.QueryRowContext(ctx, query, id).Scan(&os)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", fmt.Errorf("cannot execute query, error: %v", err)
	}

	return os, nil
}

// GetOsBySearchName Возвращать то что начинается с name или все совпадения?
func (r *Repository) GetOsBySearchName(ctx context.Context, name string) ([]os.Info, error) {
	var res []os.Info

	searchString := "%" + name + "%"

	query := "SELECT id, name FROM os WHERE name LIKE $1 LIMIT 10"

	rows, err := r.connection.QueryContext(ctx, query, searchString)

	if err != nil {
		return nil, fmt.Errorf("cannot configure query, error: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var OS os.Info
		err := rows.Scan(&OS.ID, &OS.Name)

		if err != nil {
			return nil, fmt.Errorf("cannot execute query, error: %v", err)
		}

		res = append(res, OS)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}

	return res, nil
}
