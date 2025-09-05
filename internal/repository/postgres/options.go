package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/s21platform/optionhub-service/internal/model"
)

func (r *Repository) GetOptionsByAttributeId(ctx context.Context, attributeId int64) ([]model.Option, error) {
	query, args, err := sq.Select(
		`option_id`,	
		`attribute_id`,
		`label`,
	).From(`options`).
		Where(sq.Eq{`attribute_id`: attributeId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var options []model.Option
	err = r.connection.SelectContext(ctx, &options, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch options: %w", err)
	}

	return options, nil
}