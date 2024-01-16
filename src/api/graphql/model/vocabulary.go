package model

type Vocabulary struct {
	ID uint64 `json:"id"`
	// Vocabulary's name
	Name string `json:"name"`
	// Vocabulary's title
	Title string `json:"title"`
	// Parent vocabulary
	ParentID *uint64
	// Vocabulary's description
	Description *string `json:"description"`
}

func (Vocabulary) IsEntity() {}
