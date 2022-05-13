package model

type Relation struct {
	TagId       uint64
	NamespaceId uint64
	EntityId    uint64
}

type RelationFilter struct {
	TagId     []uint64
	Namespace []string
	EntityId  []uint64
}
