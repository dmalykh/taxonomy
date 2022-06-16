package model

type Relation struct {
	Id          uint
	TagId       uint
	NamespaceId uint
	EntityId    uint
}

// EntityFilter used for external requests
// All tags specified in internal TagId's slice use "OR" operand, between TagIds "AND" operand used. See GetRelations method
type EntityFilter struct {
	TagId     [][]uint
	Namespace []string
	EntityId  []uint
	AfterId   *uint
	Limit     *uint
}

// RelationFilter used for requests to repository
type RelationFilter struct {
	TagId     [][]uint // See EntityFilter
	Namespace []uint
	EntityId  []uint
	AfterId   *uint
	Limit     *uint
}
