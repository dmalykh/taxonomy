package model

type Namespace struct {
	ID   uint64
	Data NamespaceData
}
type NamespaceData struct {
	ID    uint64
	Name  string
	Title string
}

type NamespaceFilter struct {
	Name    *string
	AfterID *uint64
	Limit   uint
}
