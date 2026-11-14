package models

type Tag struct {
	Name string
}

func NewTag(name string) *Tag {
	return &Tag{
		Name: name,
	}
}
