package model

import optionhub "github.com/s21platform/optionhub-proto/optionhub-proto"

type CategoryItem struct {
	ID    int64  `db:"id"`
	Label string `db:"name"`
}

type CategoryItemList []CategoryItem

func (c CategoryItemList) FromDTO() []*optionhub.Record {
	result := make([]*optionhub.Record, 0, len(c))

	for _, item := range c {
		result = append(result, &optionhub.Record{
			Id:    item.ID,
			Label: item.Label,
		})
	}

	return result
}
