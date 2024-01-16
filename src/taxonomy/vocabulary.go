package taxonomy

import (
	"context"
	"errors"
	"github.com/dmalykh/taxonomy/taxonomy/model"
)

var (
	ErrVocabularyNotFound   = errors.New(`vocabulary not found`)
	ErrVocabularyNotCreated = errors.New(`vocabulary had not created`)
	ErrVocabularyHasTerms   = errors.New(`vocabulary has terms, but should be empty`)
	ErrVocabularyNotUpdated = errors.New(`vocabulary had not updated`)
)

type Vocabulary interface {
	Create(ctx context.Context, data *model.VocabularyData) (*model.Vocabulary, error)
	Update(ctx context.Context, id uint64, data *model.VocabularyData) (*model.Vocabulary, error)
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*model.Vocabulary, error)
	Get(ctx context.Context, filter *model.VocabularyFilter) ([]*model.Vocabulary, error)
}
