package model

import (
	"time"

	optionhub "github.com/s21platform/optionhub-proto/optionhub/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OptionRequest struct {
	ID             int64 `db:"id"`
	AttributeID    int64 `db:"attribute_id"`
	AttributeValue string
	Value          string    `db:"value"`
	UserUuid       string    `db:"user_uuid"`
	CreatedAt      time.Time `db:"created_at"`
}

type OptionRequestList []OptionRequest

func (o OptionRequestList) ToDTO() []*optionhub.OptionRequestItem {
	result := make([]*optionhub.OptionRequestItem, 0, len(o))

	for _, item := range o {
		result = append(result, &optionhub.OptionRequestItem{
			OptionRequestId:    item.ID,
			AttributeId:        item.AttributeID,
			OptionRequestValue: item.Value,
			CreatedAt:          timestamppb.New(item.CreatedAt),
			UserUuid:           item.UserUuid,
		})
	}

	return result
}
