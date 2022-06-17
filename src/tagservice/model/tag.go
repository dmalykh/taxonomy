package model

type Tag struct {
	ID   uint
	Data TagData
}

type TagData struct {
	Name        string
	Title       string
	Description string
	CategoryID  uint
}

type TagFilter struct {
	CategoryID []uint
	Name       *string
	AfterID    *uint
	Limit      uint
	Offset     uint
}
