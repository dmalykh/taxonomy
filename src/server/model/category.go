package model

type Category struct {
	Id   uint64
	Data *CategoryData
}

type CategoryData struct {
	Name        string
	Title       string
	Description string
	PatentId    uint64
}
