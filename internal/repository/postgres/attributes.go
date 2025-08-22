package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/s21platform/optionhub-service/internal/model"
)

func (r *Repository) GetAttributesByIds(ctx context.Context, attributesIds []int64) ([]model.Attribute, error) {
	query, args, err := sq.Select(
		`attribute_id`,
		`name`,
		`type`,
		`description`,
		`allowed_operators`,
	).From(`attributes`).
		Where(sq.Eq{`attribute_id`: attributesIds}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var attributesMeta []model.Attribute
	err = r.connection.SelectContext(ctx, &attributesMeta, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch attributes meta: %w", err)
	}
	return attributesMeta, nil
}

func (r *Repository) GetEntityAttributesByIds(ctx context.Context, attributesEntityIds []int64) ([]model.EntityAttribute, error) {
	if len(attributesEntityIds) == 0 {
		return nil, nil
	}

	query, args, err := sq.Select(
		`entity_attribute_id`,
		`entity_type`,
		`attribute_id`,
		`label`,
		`is_required`,
		`order_index`,
		`visibility_rules`,
	).From(`entity_attributes`).
		Where(sq.Eq{`entity_attribute_id`: attributesEntityIds}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}
	var entityAttributes []model.EntityAttribute
	err = r.connection.SelectContext(ctx, &entityAttributes, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch entity attributes: %w", err)
	}
	return entityAttributes, nil
}
