package repository

import (
	"context"
	"errors"
	"github.com/dmalykh/taxonomy/taxonomy/model"
)

var (
	ErrCreateTerm = errors.New(`failed to create term`)
	ErrUpdateTerm = errors.New(`failed to update term`)
	ErrFindTerm   = errors.New(`failed to find term`)
	ErrDeleteTerm = errors.New(`failed to delete term`)
)

type Term interface {
	Create(ctx context.Context, data *model.TermData) (*model.Term, error)
	Update(ctx context.Context, id uint64, data *model.TermData) (*model.Term, error)
	Delete(ctx context.Context, filter *TermFilter) error
	Get(ctx context.Context, filter *TermFilter) ([]*model.Term, error)
}

type TermFilter struct {
	ID           []uint64 // anyOf
	VocabularyID []uint64 // anyOf
	SuperID      []uint64 // anyOf
	SubID        []uint64 // anyOf
	Name         *string
	AfterID      *uint64
	Limit        uint
	Offset       uint
}
