package model

type Term struct {
	ID uint64 `json:"id"`
	// Term's name
	Name string `json:"name"`
	// Term's title
	Title *string `json:"title"`
	// Term's vocabulary
	VocabularyID uint64 `json:"vocabulary"`
	// Description
	Description *string `json:"description"`
}

func (t Term) IsEntity() {}
