package postgres

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL

	"github.com/s21platform/optionhub-service/internal/config"
	"github.com/s21platform/optionhub-service/internal/model"
)

const (
	attributeValuesTable = "attribute_values"
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

func (r *Repository) AddAttributeValue(ctx context.Context, in model.AttributeValue) error {
	queryTmp := sq.Insert("attribute_values").
		Columns("attribute_id", "value").
		Values(in.AttributeId, in.Value)

	if in.ParentId != nil {
		queryTmp = queryTmp.Columns("parent_id").Values(*in.ParentId)
	}

	sqlQuery, args, err := queryTmp.PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return fmt.Errorf("failed to build SQL query: %v", err)
	}

	_, err = r.connection.ExecContext(ctx, sqlQuery, args...)

	if err != nil {
		return fmt.Errorf("failed to add attribute into postgres: %v", err)
	}

	return nil
}

func (r *Repository) GetValuesByAttributeId(ctx context.Context, attributeId int64) (model.AttributeValueList, error) {
	var values model.AttributeValueList

	query, args, err := sq.
		Select(
			"id",
			"value",
			"parent_id",
		).
		From(attributeValuesTable).
		Where(sq.Eq{"attribute_id": attributeId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	err = r.connection.SelectContext(ctx, &values, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	return values, nil
}
