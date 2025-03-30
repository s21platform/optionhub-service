package model

type Attribute struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}
