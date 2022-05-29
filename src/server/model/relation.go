package model

type Relation struct {
	TagId       uint
	NamespaceId uint
	EntityId    uint
}

type RelationFilter struct {
	TagId     []uint
	Namespace []string
	EntityId  []uint
}
