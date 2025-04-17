package model

import (
	"github.com/samber/lo"

	optionhubproto_v1 "github.com/s21platform/optionhub-proto/optionhub/v1"
)

type Attribute struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type AttributeValue struct {
	Id          int64  `db:"id"`
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

type AttributeValueList []AttributeValue

func (a AttributeValueList) FromDTO() []*optionhubproto_v1.Option {
	result := make([]*optionhubproto_v1.Option, 0)

	roots := lo.Filter(a, func(val AttributeValue, _ int) bool {
		return val.ParentId == nil
	})

	//group children by parent_id
	childrenMap := make(map[int64]AttributeValueList)
	for _, val := range a {
		if val.ParentId != nil {
			childrenMap[*val.ParentId] = append(childrenMap[*val.ParentId], val)
		}
	}

	for _, root := range roots {
		rootNode := optionhubproto_v1.Option{
			OptionId:    root.Id,
			OptionValue: root.Value,
			Children:    buildTree(root.Id, childrenMap),
		}
		result = append(result, &rootNode)
	}

	return result
}

func buildTree(parentId int64, children map[int64]AttributeValueList) []*optionhubproto_v1.Option {
	result := make([]*optionhubproto_v1.Option, 0)
	for _, child := range children[parentId] {
		node := optionhubproto_v1.Option{
			OptionId:    child.Id,
			OptionValue: child.Value,
			Children:    buildTree(child.Id, children),
		}
		result = append(result, &node)
	}
	return result
}
