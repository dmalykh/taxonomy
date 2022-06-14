package model

type Category struct {
	Id   uint
	Data CategoryData
}

type CategoryData struct {
	Name        string
	Title       string
	Description *string
	ParentId    *uint
}
