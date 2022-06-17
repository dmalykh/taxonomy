package model

type Category struct {
	ID int64 `json:"id"`
	// Category's name
	Name string `json:"name"`
	// Category's title
	Title string `json:"title"`
	// Parent category
	ParentID *int64
	// Category's description
	Description *string `json:"description"`
}

func (Category) IsEntity() {}
