package model

import optionhubproto_v1 "github.com/s21platform/optionhub-proto/optionhub/v1"

type Attribute struct {
	AttributeId int64  `db:"attribute_id"`
	Value       string `db:"value"`
	ParentId    int64  `db:"parent_id"`
}

type SetAttributeMessage struct {
	AttributeId int64 `json:"attribute_id"`
}

func (a *Attribute) AttributeToDTO(in *optionhubproto_v1.SetAttributeByIdIn) (Attribute, error) {
	result := Attribute{
		AttributeId: in.AttributeId,
		Value:       in.Value,
		ParentId:    in.ParentId,
	}
	return result, nil
}
