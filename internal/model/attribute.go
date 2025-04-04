package model

import (
	optionhubproto_v1 "github.com/s21platform/optionhub-proto/optionhub/v1"
)

type Attribute struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type AttributeValue struct {
	AttributeId int64  `db:"attribute_id"`
	Value       string `db:"value"`
	ParentId    *int64 `db:"parent_id"`
}

func (a *AttributeValue) ToDTO(in *optionhubproto_v1.AddAttributeValueIn) (AttributeValue, error) {
	result := AttributeValue{
		AttributeId: in.AttributeId,
		Value:       in.Value,
		ParentId:    in.ParentId,
	}
	return result, nil
}
