package repository

import (
	"context"
	"errors"
	"github.com/dmalykh/taxonomy/taxonomy/model"
)

var (
	ErrCreateReference  = errors.New(`failed to create reference`)
	ErrGetReference     = errors.New(`failed to get reference`)
	ErrWithoutNamespace = errors.New(`namespace required`)
	ErrDeleteReferences = errors.New(`failed to delete reference`)
)

type Reference interface {
	Create(ctx context.Context, reference ...*ReferenceModel) error
	Delete(ctx context.Context, filter *ReferenceFilter) error
	Get(ctx context.Context, filter *ReferenceFilter) ([]ReferenceModel, error)
}

// ReferenceFilter used for requests to repository.
// All terms specified in internal TermID's slice use "OR" operand, between TermIDs "AND" operand used.
type ReferenceFilter struct {
	TermID      [][]uint64 // {{X OR X} AND {X OR X OR X}}
	NamespaceID []uint64   // Required!
	EntityID    []model.EntityID
	AfterID     *uint64
	Limit       uint
}

type ReferenceModel struct {
	ID          uint64
	TermID      uint64
	NamespaceID uint64
	EntityID    model.EntityID
}
