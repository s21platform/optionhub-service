package model

type Option struct {
	OptionID int64 `db:"option_id"`
	AttributeID int64 `db:"attribute_id"`
	Label string `db:"label"`
}