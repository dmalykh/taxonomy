package model

type Tag struct {
	Id   uint
	Data TagData
}

type TagData struct {
	Name        string
	Title       string
	Description string
	CategoryId  uint
}

type TagFilter struct {
	CategoryId []uint
	Name       *string
	AfterId    *uint
	Limit      *uint
	Offset     *uint
}
