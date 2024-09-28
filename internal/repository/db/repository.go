package db

import (
	"context"
	"database/sql"
	"fmt"
	optionhub_proto "github.com/s21platform/optionhub-proto/optionhub-proto"
	"log"
	"time"

	_ "github.com/lib/pq"
	"optionhub-service/internal/config"
)

type Repository struct {
	сonnection *sql.DB
}

func connect(cfg *config.Config) (*Repository, error) {
	//Connect db
	conStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database, cfg.Postgres.Host, cfg.Postgres.Port)

	db, err := sql.Open("postgres", conStr)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	//Сhecking connection db
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping: %w", err)
	}
	return &Repository{db}, nil
}

func (r *Repository) Close() {
	r.сonnection.Close()
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

// пока без user_uuid (будет браться из токена)
func (r *Repository) AddOS(ctx context.Context, name string) (int64, error) {
	query := "INSERT INTO os(name, create_at) VALUES ($1, $2) RETURNING id"
	createTime := time.Now().UTC()
	var id int64

	err := r.сonnection.QueryRowContext(ctx, query, name, createTime).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("cannot execute query, error: %v", err)
	}
	return id, nil
}

func (r *Repository) GetOsById(ctx context.Context, id int64) (string, error) {
	var os string
	query := "SELECT name from os where id = $1"

	err := r.сonnection.QueryRowContext(ctx, query, id).Scan(&os)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("cannot execute query, error: %v", err)
	}
	return os, nil
}

// лимит можно передавать как параметр позже
// нужно возвращать то что начинается на name или все строки с совпадением?
func (r *Repository) GetOsBSearchName(ctx context.Context, name string) ([]*optionhub_proto.Record, error) {
	var records []*optionhub_proto.Record
	searchString := "%" + name + "%"

	query := "SELECT id, name FROM os WHERE name LIKE $1 LIMIT 10"
	rows, err := r.сonnection.QueryContext(ctx, query, searchString)
	if err != nil {
		return nil, fmt.Errorf("cannot configure query, error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		record := optionhub_proto.Record{}
		err := rows.Scan(&record.Id, &record.Value)
		if err != nil {
			return nil, fmt.Errorf("cannot execute query, error: %v", err)
		}
		records = append(records, &record)
		fmt.Println(&record)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return records, nil
}
