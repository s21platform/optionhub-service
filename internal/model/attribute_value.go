package model

import (
	"github.com/samber/lo"

	optionhub "github.com/s21platform/optionhub-proto/optionhub/v1"
)

type AttributeValue struct {
	Id       int64  `db:"id"`
	Value    string `db:"value"`
	ParentId *int64 `db:"parent_id"` //если нет родителя то 0 или нил?
	//нужно ли еще created_at селектить?
}

type AttributeValueList []AttributeValue

func (a AttributeValueList) ToDTO() []*optionhub.Option {
	result := make([]*optionhub.Option, 0)

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
		rootNode := optionhub.Option{
			OptionId:    root.Id,
			OptionValue: root.Value,
			Children:    buildTree(root.Id, childrenMap),
		}
		result = append(result, &rootNode)
	}

	return result
}

func buildTree(parentId int64, children map[int64]AttributeValueList) []*optionhub.Option {
	result := make([]*optionhub.Option, 0)
	for _, child := range children[parentId] {
		node := optionhub.Option{
			OptionId:    child.Id,
			OptionValue: child.Value,
			Children:    buildTree(child.Id, children),
		}
		result = append(result, &node)
	}
	return result
}
