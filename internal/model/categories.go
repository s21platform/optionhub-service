package model

import optionhub "github.com/s21platform/optionhub-proto/optionhub-proto"

type CategoryItem struct {
	ID    int64  `db:"id"`
	Label string `db:"name"`
}

type CategoryItemList []CategoryItem

func (c CategoryItemList) FromDTO() []*optionhub.Record {
	result := make([]*optionhub.Record, 0, len(c))

	for _, avatar := range c {
		result = append(result, avatar.ToProto())
	}

	return result
}

func (c *CategoryItem) ToProto() *optionhub.Record {
	return &optionhub.Record{
		Id:    c.ID,
		Label: c.Label,
	}
}
