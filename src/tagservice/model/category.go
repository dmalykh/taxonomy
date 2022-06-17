package model

type Category struct {
	ID   uint
	Data CategoryData
}

type CategoryData struct {
	Name        string
	Title       string
	Description *string
	ParentID    *uint
}

type CategoryFilter struct {
	ParentID *uint
	Name     *string
}
