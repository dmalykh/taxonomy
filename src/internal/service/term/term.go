package term

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
	ReferenceService  taxonomy.Reference
	TermRepository    repository.Term
	VocabularyService taxonomy.Vocabulary
	Logger            *zap.Logger
}

func New(config *Config) taxonomy.Term {
	return &TermService{
		referenceService:  config.ReferenceService,
		vocabularyService: config.VocabularyService,
		termRepository:    config.TermRepository,
		log:               config.Logger,
	}
}

type TermService struct {
	referenceService  taxonomy.Reference
	vocabularyService taxonomy.Vocabulary
	namespaceService  taxonomy.Namespace
	termRepository    repository.Term
	log               *zap.Logger
}

func (t *TermService) Create(ctx context.Context, data *model.TermData) (model.Term, error) {
	logger := t.log.With(zap.String(`method`, `Create`), zap.Any(`data`, *data))

	if err := t.checkVocabularies(ctx, data.VocabularyID); err != nil {
		return model.Term{}, err
	}

	term, err := t.termRepository.Create(ctx, data)
	logger.Debug(`term created`, zap.Any(`term`, term), zap.Error(err))

	if err != nil {
		return model.Term{}, fmt.Errorf(`%w %s`, taxonomy.ErrTermNotCreated, err.Error())
	}

	return term, nil
}

func (t *TermService) checkVocabularies(ctx context.Context, vocabulariesID []uint64) error {
	// Check vocabularies exists
	for _, id := range vocabulariesID {
		if _, err := t.vocabularyService.GetByID(ctx, id); err != nil {
			if errors.Is(err, repository.ErrFindVocabulary) {
				return fmt.Errorf(`%w %d`, taxonomy.ErrVocabularyNotFound, id)
			}

			return fmt.Errorf(`unknown vocabulary error %w`, err)
		}
	}

	return nil
}

func (t *TermService) Update(ctx context.Context, id uint64, data *model.TermData) (*model.Term, error) {
	logger := t.log.With(zap.String(`method`, `Update`), zap.Uint64("id", id),
		zap.Any(`data`, *data))

	// Check term exists
	terms, err := t.termRepository.Get(ctx, &repository.TermFilter{ID: []uint64{id}})
	if err != nil {
		logger.Error(`get term by id`, zap.Error(err))

		if errors.Is(err, repository.ErrFindTerm) {
			return nil, fmt.Errorf(`%w %d`, taxonomy.ErrTermNotFound, id)
		}

		return nil, fmt.Errorf(`unknown error %w`, err)
	}

	if len(terms) != 1 {
		return nil, fmt.Errorf(`%w, got %d results`, taxonomy.ErrTermNotFound, len(terms))
	}

	var term = terms[0]

	// Avoid empty values
	if data.Name == `` {
		data.Name = term.Data.Name
	}

	if data.Title == `` {
		data.Title = term.Data.Title
	}

	if data.Description == `` {
		data.Description = term.Data.Description
	}

	if len(data.VocabularyID) == 0 {
		data.VocabularyID = term.Data.VocabularyID
	} else {
		// Check vocabulary exists
		if err := t.checkVocabularies(ctx, data.VocabularyID); err != nil {
			return nil, err
		}
	}

	// Update term
	updated, err := t.termRepository.Update(ctx, term.ID, data)
	logger.Debug(`term updated`, zap.Any(`term`, updated), zap.Error(err))

	if err != nil {
		return nil, errors.Join(taxonomy.ErrTermNotUpdated, err)
	}

	return updated, nil
}

func (t *TermService) Delete(ctx context.Context, id uint64) error {
	logger := t.log.With(zap.String(`method`, `Delete`), zap.Uint64("id", id))

	// Check term exists
	term, err := t.termRepository.GetByID(ctx, id)
	if err != nil {
		logger.Error(`get term by id`, zap.Error(err))

		if errors.Is(err, repository.ErrFindTerm) {
			return fmt.Errorf(`%w %d`, taxonomy.ErrTermNotFound, id)
		}

		return fmt.Errorf(`unknown error %w`, err)
	}

	// Reference exists check
	logger.Debug(`check references`, zap.Uint64(`id`, term.ID))

	ref, err := t.referenceService.Get(ctx, &model.ReferenceFilter{TermID: [][]uint64{{term.ID}}})
	if err != nil {
		logger.Error(`get references by term_id`, zap.Uint64(`term_id`, term.ID), zap.Error(err))

		return fmt.Errorf(`get references by term error: %w`, err)
	}

	if len(ref) > 0 {
		return fmt.Errorf(`can't remove term %q: %d %w`, term.ID, len(ref), taxonomy.ErrTermReferenceExists)
	}

	// Delete term
	logger.Debug(`delete term by id`, zap.Uint64(`id`, term.ID))

	if err := t.termRepository.DeleteByID(ctx, term.ID); err != nil {
		logger.Error(`delete term by id error`, zap.Uint64(`term_id`, term.ID), zap.Error(err))

		return fmt.Errorf(`can't remove term %w`, err)
	}

	return nil
}

func (t *TermService) GetByID(ctx context.Context, id uint64) (model.Term, error) {
	logger := t.log.With(zap.String(`method`, `GetByID`), zap.Uint64("id", id))

	term, err := t.termRepository.GetByID(ctx, id)
	if err != nil {
		logger.Error(`get term by id`, zap.Error(err))

		if errors.Is(err, repository.ErrFindTerm) {
			return term, fmt.Errorf(`%w %d`, taxonomy.ErrTermNotFound, id)
		}

		return term, fmt.Errorf(`unknown error %w`, err)
	}

	return term, nil
}

func (t *TermService) Get(ctx context.Context, filter *model.TermFilter) ([]model.Term, error) {
	logger := t.log.With(zap.String(`method`, `Get`), zap.Any(`filter`, filter))

	terms, err := t.termRepository.GetList(ctx, filter)
	logger.Debug(`got terms`, zap.Any(`terms`, terms), zap.Error(err))

	if err != nil {
		return nil, fmt.Errorf(`unknown error %w`, err)
	}

	return terms, nil
}
