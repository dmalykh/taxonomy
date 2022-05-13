package model

type Tag struct {
	Id   uint64
	Data *TagData
}

type TagData struct {
	Name        string
	Title       string
	Description string
	CategoryId  uint64
	Active      bool
}

type TagFilter struct {
	CategoryId []uint64
	Active     *bool
	Limit      uint64
	Offset     uint64
}
