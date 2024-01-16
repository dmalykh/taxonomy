package vocabulary

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/taxonomy/taxonomy"
	"github.com/dmalykh/taxonomy/taxonomy/model"
	"github.com/dmalykh/taxonomy/taxonomy/repository"

	"go.uber.org/zap"
)

type Config struct {
	VocabularyRepository repository.Vocabulary
	TermService          taxonomy.Term
	Logger               *zap.Logger
}

func New(config *Config) taxonomy.Vocabulary {
	return &VocabularyService{
		termService:          config.TermService,
		vocabularyRepository: config.VocabularyRepository,
		log:                  config.Logger,
	}
}

type VocabularyService struct {
	termService          taxonomy.Term
	vocabularyRepository repository.Vocabulary
	log                  *zap.Logger
}

func (c *VocabularyService) Create(ctx context.Context, data *model.VocabularyData) (*model.Vocabulary, error) {
	logger := c.log.With(zap.String(`method`, `Create`), zap.Any(`data`, *data))
	// Check parent's vocabulary exists
	if data.ParentID != nil {
		if ok, err := c.exists(ctx, *data.ParentID); !ok || err != nil {
			if err != nil {
				return nil, fmt.Errorf(`%w, get parent id  %d error: %w`,
					taxonomy.ErrVocabularyNotFound, *data.ParentID, err)
			}

			return nil, fmt.Errorf(`id %d %w`, *data.ParentID, taxonomy.ErrVocabularyNotFound)
		}
	}
	// Create vocabulary
	vocabulary, err := c.vocabularyRepository.Create(ctx, data)
	logger.Debug(`vocabulary created`, zap.Error(err))

	if err != nil {
		return nil, fmt.Errorf(`%w %w`, taxonomy.ErrVocabularyNotCreated, err)
	}

	return vocabulary, nil
}

func (c *VocabularyService) Update(ctx context.Context, id uint64, data *model.VocabularyData) (*model.Vocabulary, error) {
	logger := c.log.With(zap.String(`method`, `Update`), zap.Uint64("id", id))

	vocabularies, err := c.vocabularyRepository.Get(ctx, &repository.VocabularyFilter{ID: []uint64{id}})
	if err != nil {
		logger.Error(`get vocabulary by id`, zap.Error(err))

		if errors.Is(err, repository.ErrFindVocabulary) {
			return nil, fmt.Errorf(`%w %d`, taxonomy.ErrVocabularyNotFound, id)
		}

		return nil, fmt.Errorf(`unknown error %w`, err)
	}

	if len(vocabularies) != 1 {
		return nil, fmt.Errorf(`%w, got %d results`, taxonomy.ErrVocabularyNotFound, len(vocabularies))
	}

	var vocabulary = vocabularies[0]

	// Check parent's vocabulary exists
	if data.ParentID != nil {
		if ok, err := c.exists(ctx, *data.ParentID); !ok || err != nil {
			if err != nil {
				return nil, fmt.Errorf(`%w, get parent id  %d error: %w`,
					taxonomy.ErrVocabularyNotFound, *data.ParentID, err)
			}

			return nil, fmt.Errorf(`id %d %w`, *data.ParentID, taxonomy.ErrVocabularyNotFound)
		}
	}

	// Avoid empty values
	if data.Name == `` {
		data.Name = vocabulary.Data.Name
	}

	if data.Title == `` {
		data.Title = vocabulary.Data.Title
	}

	if data.Description == nil {
		data.Description = vocabulary.Data.Description
	}

	if data.ParentID == nil {
		data.ParentID = vocabulary.Data.ParentID
	} else if *data.ParentID == 0 {
		data.ParentID = nil
	}

	// Avoid loops with ParentID
	if data.ParentID != nil && *data.ParentID == id {
		return nil, fmt.Errorf(`parentid (%d) can't equals id (%d)`, *data.ParentID, id)
	}

	vocabulary, err = c.vocabularyRepository.Update(ctx, vocabulary.ID, data)
	logger.Debug(`vocabulary updated`, zap.Error(err))

	if err != nil {
		return nil, fmt.Errorf(`%w %w`, taxonomy.ErrVocabularyNotUpdated, err)
	}

	return vocabulary, nil
}

// Delete vocabulary and it's dependencies.
func (c *VocabularyService) Delete(ctx context.Context, id uint64) error {
	logger := c.log.With(zap.String(`method`, `Delete`), zap.Uint64("id", id))
	// Check vocabulary exists
	if _, err := c.GetByID(ctx, id); err != nil {
		return err
	}

	// Check terms. Vocabulary should be empty before deletion
	terms, err := c.termService.Get(ctx, &model.TermFilter{VocabularyID: []uint64{id}})
	logger.Debug(`get terms of vocabulary`, zap.Error(err))

	if err != nil {
		return fmt.Errorf(`unknown error %w`, err)
	}

	if len(terms) > 0 {
		return taxonomy.ErrVocabularyHasTerms
	}

	// Delete vocabulary
	logger.Debug(`delete vocabulary`, zap.Uint64(`id`, id))

	if err := c.vocabularyRepository.Delete(ctx, &repository.VocabularyFilter{ID: []uint64{id}}); err != nil {
		return fmt.Errorf(`can't remove vocabulary %w`, err)
	}

	return nil
}

func (c *VocabularyService) Get(ctx context.Context, filter *model.VocabularyFilter) ([]*model.Vocabulary, error) {
	logger := c.log.With(zap.String(`method`, `Get`), zap.Any(`filter`, filter))

	list, err := c.vocabularyRepository.Get(ctx, &repository.VocabularyFilter{
		Name:     valToSlice[string](filter.Name),
		ParentID: valToSlice[uint64](filter.ParentID),
	})
	logger.Debug(`get list`, zap.Error(err))

	if err != nil {
		return nil, fmt.Errorf(`can't receive list of vocabularys %w`, err)
	}

	return list, nil
}

func (c *VocabularyService) GetByID(ctx context.Context, id uint64) (*model.Vocabulary, error) {
	logger := c.log.With(zap.String(`method`, `GetByID`), zap.Uint64("id", id))

	vocabularies, err := c.vocabularyRepository.Get(ctx, &repository.VocabularyFilter{ID: []uint64{id}})
	if err != nil {
		logger.Error(`get vocabulary by id`, zap.Error(err))

		if errors.Is(err, repository.ErrFindVocabulary) {
			return nil, fmt.Errorf(`%w %d`, taxonomy.ErrVocabularyNotFound, id)
		}

		return nil, fmt.Errorf(`unknown error %w`, err)
	}

	if len(vocabularies) != 1 {
		return nil, fmt.Errorf(`%w, got %d results`, taxonomy.ErrVocabularyNotFound, len(vocabularies))
	}

	return nil, nil
}

func (c *VocabularyService) exists(ctx context.Context, id uint64) (bool, error) {
	vocabularies, err := c.vocabularyRepository.Get(ctx, &repository.VocabularyFilter{ID: []uint64{id}})
	if err != nil {
		if errors.Is(err, repository.ErrFindVocabulary) {
			return false, nil
		}

		return false, fmt.Errorf(`unknown parent id error %w`, err)
	}

	if len(vocabularies) != 1 {
		return false, fmt.Errorf(`%w, got %d results`, taxonomy.ErrVocabularyNotFound, len(vocabularies))
	}

	return true, nil
}

func valToSlice[T any](input *T) []T {
	if input != nil {
		return []T{*input}
	}
	return []T{}
}
