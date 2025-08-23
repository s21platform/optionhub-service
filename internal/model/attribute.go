package model

type EntityAttribute struct {
	EntityAttributeID int64  `db:"entity_attribute_id"`
	EntityType        string `db:"entity_type"`
	AttributeID       int64  `db:"attribute_id"`
	Label             string `db:"label"`
	IsRequired        bool   `db:"is_required"`
	OrderIndex        int    `db:"order_index"`
	VisibilityRules   []byte `db:"visibility_rules"`
}

type Attribute struct {
	ID               int64                    `db:"attribute_id"`
	Name             string                   `db:"name"`
	Type             string                   `db:"type"`
	Description      *string                  `db:"description"`
	AllowedOperators *AttributeAllowOperators `db:"allowed_operators"`
}

type AttributeAllowOperators struct {
	// заглушка
	Stub bool `json:"stub"`
}
