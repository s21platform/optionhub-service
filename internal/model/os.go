package model

import optionhub "github.com/s21platform/optionhub-proto/optionhub-proto"

type OS struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type OSList []OS

func (o *OSList) FromDTO() []*optionhub.Record {
	result := make([]*optionhub.Record, 0, len(*o))

	for _, avatar := range *o {
		result = append(result, avatar.ToProto())
	}

	return result
}

func (os *OS) ToProto() *optionhub.Record {
	return &optionhub.Record{
		Id:   os.ID,
		Name: os.Name,
	}
}
