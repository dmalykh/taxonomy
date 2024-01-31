package taxonomy

import (
	"context"
	"errors"
	"github.com/dmalykh/taxonomy/taxonomy/model"
)

var (
	ErrTermNotFound   = errors.New(`term not found`)
	ErrTermNotCreated = errors.New(`term had not created`)
	ErrTermNotUpdated = errors.New(`term have not updated`)
)

type Term interface {
	Create(ctx context.Context, data *model.TermData) (*model.Term, error)
	Update(ctx context.Context, id uint64, data *model.TermData) (*model.Term, error)
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*model.Term, error)

	// Get returns slice with terms that proper for conditions. Set nil vocabulary_id to receive terms from all categories.
	Get(ctx context.Context, filter *model.TermFilter) ([]*model.Term, error)
}
