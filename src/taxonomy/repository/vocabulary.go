package repository

import (
	"context"
	"errors"
	"github.com/dmalykh/taxonomy/taxonomy/model"
)

var (
	ErrCreateVocabulary = errors.New(`failed to create vocabulary`)
	ErrUpdateVocabulary = errors.New(`failed to update vocabulary`)
	ErrNotUniqueName    = errors.New(`vocabulary's name and parent must be unique`)
	ErrFindVocabulary   = errors.New(`failed to find vocabulary`)
	ErrDeleteVocabulary = errors.New(`failed to delete vocabulary`)
)

type Vocabulary interface {
	Create(ctx context.Context, data *model.VocabularyData) (*model.Vocabulary, error)
	Update(ctx context.Context, id uint64, data *model.VocabularyData) (*model.Vocabulary, error)
	Delete(ctx context.Context, filter *VocabularyFilter) error
	Get(ctx context.Context, filter *VocabularyFilter) ([]*model.Vocabulary, error)
}

type VocabularyFilter struct {
	ID       []uint64
	ParentID []uint64
	Name     []string
}
