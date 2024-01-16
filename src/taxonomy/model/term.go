package model

type Term struct {
	ID   uint64
	Data TermData
}

type TermData struct {
	Name         string
	Title        string
	Description  string
	VocabularyID []uint64
	SuperID      []uint64
	SubID        []uint64
}

type TermFilter struct {
	VocabularyID []uint64 // anyOf
	SuperID      []uint64 // anyOf
	SubID        []uint64 // anyOf
	Name         *string
	AfterID      *uint64
	Limit        uint
	Offset       uint
}
