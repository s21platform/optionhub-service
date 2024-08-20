package db

import (
	"database/sql"
	"fmt"
	uuid "github.com/satori/go.uuid"
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

func (r *Repository) AddedOS(id int, name string, userUID uuid.UUID) error {
	createTime := time.Now()

	_, err := r.сonnection.Exec("INSERT INTO os(id, name, user_uuid, create_at) VALUES ($1, $2, $3 $4)", id, name, userUID, createTime)
	if err != nil {
		return err
	}
	return nil
}
