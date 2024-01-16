package model

type Vocabulary struct {
	ID   uint64
	Data VocabularyData
}

type VocabularyData struct {
	Name        string
	Title       string
	Description *string
	ParentID    *uint64
}

type VocabularyFilter struct {
	ParentID *uint64
	Name     *string
}
