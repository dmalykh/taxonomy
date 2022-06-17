package model

type Relation struct {
	ID          uint
	TagID       uint
	NamespaceID uint
	EntityID    uint
}

// EntityFilter used for external requests
// All tags specified in internal TagID's slice use "OR" operand, between TagIds "AND" operand used. See GetRelations method.
type EntityFilter struct {
	TagID     [][]uint
	Namespace []string
	EntityID  []uint
	AfterID   *uint
	Limit     *uint
}

// RelationFilter used for requests to repository.
type RelationFilter struct {
	TagID     [][]uint // See EntityFilter
	Namespace []uint
	EntityID  []uint
	AfterID   *uint
	Limit     *uint
}
