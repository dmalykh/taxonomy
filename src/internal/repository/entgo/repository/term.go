package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/taxonomy/internal/repository/entgo/ent"
	"github.com/dmalykh/taxonomy/internal/repository/entgo/ent/predicate"
	"github.com/dmalykh/taxonomy/internal/repository/entgo/ent/term"
	"github.com/dmalykh/taxonomy/internal/repository/entgo/ent/vocabulary"
	"github.com/dmalykh/taxonomy/taxonomy/model"
	"github.com/dmalykh/taxonomy/taxonomy/repository"
)

type Term struct {
	client *ent.TermClient
}

func NewTerm(client *ent.TermClient) repository.Term {
	return &Term{
		client: client,
	}
}

func (t *Term) Create(ctx context.Context, data *model.TermData) (*model.Term, error) {
	created, err := t.client.Create().
		SetName(data.Name).
		SetTitle(data.Title).
		SetDescription(data.Description).
		AddVocabularyIDs(data.VocabularyID...).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", repository.ErrCreateTerm, err.Error())
	}

	trm, err := t.client.Query().WithVocabulary().Where(term.ID(created.ID)).Only(ctx)
	if err != nil {
		return nil, errors.Join(repository.ErrFindTerm, err)
	}

	return t.ent2model(trm), nil
}

func (t *Term) Update(ctx context.Context, id uint64, data *model.TermData) (*model.Term, error) {
	updated, err := t.client.UpdateOneID(id).
		SetName(data.Name).
		SetTitle(data.Title).
		SetDescription(data.Description).
		AddVocabularyIDs(data.VocabularyID...).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", repository.ErrUpdateTerm, err.Error())
	}

	trm, err := t.client.Query().WithVocabulary().Where(term.ID(updated.ID)).Only(ctx)
	if err != nil {
		return nil, errors.Join(repository.ErrFindTerm, err)
	}

	return t.ent2model(trm), err
}

func (t *Term) Delete(ctx context.Context, filter *repository.TermFilter) error {
	_, err := t.client.Delete().Where(
		t.buildQuery(filter)...,
	).Exec(ctx)
	if err != nil {
		return errors.Join(repository.ErrDeleteTerm, err)
	}

	return nil
}

func (t *Term) Get(ctx context.Context, filter *repository.TermFilter) ([]*model.Term, error) {
	entterms, err := t.client.Query().Where(
		t.buildQuery(filter)...,
	).
		WithVocabulary().
		All(ctx)
	if err != nil {
		return nil, errors.Join(repository.ErrFindTerm, err)
	}

	terms := make([]*model.Term, 0, len(entterms))

	for _, entterm := range entterms {
		terms = append(terms, t.ent2model(entterm))
	}

	return terms, nil
}

func (t *Term) buildQuery(filter *repository.TermFilter) []predicate.Term {
	var predicates = make([]predicate.Term, 0)
	// Filter by id
	if len(filter.ID) > 0 {
		predicates = append(predicates, term.IDIn(filter.ID...))
	}
	// Filter by vocabulary id
	if len(filter.VocabularyID) > 0 {
		predicates = append(predicates, term.HasVocabularyWith(
			vocabulary.IDIn(filter.VocabularyID...),
		))
	}
	// Filter by name
	if filter.Name != nil {
		predicates = append(predicates, term.Name(*filter.Name))
	}
	// Get subterms that have certain super
	if len(filter.SuperID) > 0 {
		predicates = append(predicates, term.HasSupertermsWith(
			term.IDIn(filter.SuperID...),
		))
	}
	// Get superterms for certain term
	if len(filter.SubID) > 0 {
		predicates = append(predicates, term.HasSubtermsWith(
			term.IDIn(filter.SubID...),
		))
	}
	// Add condition to id
	if filter.AfterID != nil {
		predicates = append(predicates, term.IDGT(*filter.AfterID))
	}

	return predicates
}

func (t *Term) ent2model(term *ent.Term) *model.Term {
	return &model.Term{
		ID: term.ID,
		Data: model.TermData{
			Name:        term.Name,
			Title:       term.Title,
			Description: term.Description,
			VocabularyID: toUint64s[ent.Vocabulary](term.Edges.Vocabulary, func(item *ent.Vocabulary) uint64 {
				return item.ID
			}),
			SuperID: toUint64s[ent.Term](term.Edges.Superterms, func(item *ent.Term) uint64 {
				return item.ID
			}),
			SubID: toUint64s[ent.Term](term.Edges.Subterms, func(item *ent.Term) uint64 {
				return item.ID
			}),
		},
	}
}
