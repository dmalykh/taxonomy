package model

type Tag struct {
	ID int64 `json:"id"`
	// Tag's name
	Name string `json:"name"`
	// Tag's title
	Title *string `json:"title"`
	// Tag's category
	CategoryID int64 `json:"category"`
	// Description
	Description *string `json:"description"`
}

func (t Tag) IsEntity() {}
