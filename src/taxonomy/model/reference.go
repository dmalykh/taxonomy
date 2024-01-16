package model

type Reference struct {
	ID        uint64
	TermID    uint64
	Namespace string
	EntityID  EntityID
}

// ReferenceFilter used for requests to repository.
// All terms specified in internal TermID's slice use "OR" operand, between TermIDs "AND" operand used. See GetReferences method.
type ReferenceFilter struct {
	TermID    [][]uint64 // See EntityFilter
	Namespace []string   // OR operand if used
	EntityID  []EntityID // OR operand if used
	AfterID   *uint64
	Limit     *uint
}
