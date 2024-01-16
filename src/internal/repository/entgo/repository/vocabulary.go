package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/taxonomy/internal/repository/entgo/ent"
	"github.com/dmalykh/taxonomy/internal/repository/entgo/ent/predicate"
	"github.com/dmalykh/taxonomy/internal/repository/entgo/ent/vocabulary"
	"github.com/dmalykh/taxonomy/taxonomy/model"
	"github.com/dmalykh/taxonomy/taxonomy/repository"
)

type Vocabulary struct {
	client *ent.VocabularyClient
}

func NewVocabulary(client *ent.VocabularyClient) repository.Vocabulary {
	return &Vocabulary{
		client: client,
	}
}

func (v *Vocabulary) Create(ctx context.Context, data *model.VocabularyData) (*model.Vocabulary, error) {
	ns, err := v.client.Create().
		SetName(data.Name).
		SetTitle(data.Title).
		SetNillableDescription(data.Description).
		SetNillableParentID(func() *uint64 { return data.ParentID }()).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, repository.ErrNotUniqueName
		}

		return nil, fmt.Errorf("%w: %s", repository.ErrCreateVocabulary, err.Error())
	}

	return v.ent2model(ns), nil
}

func (v *Vocabulary) Update(ctx context.Context, id uint64, data *model.VocabularyData) (*model.Vocabulary, error) {
	updated, err := v.client.UpdateOneID(id).
		SetName(data.Name).
		SetTitle(data.Title).
		SetNillableDescription(data.Description).
		ClearParentID().
		SetNillableParentID(func() *uint64 { return data.ParentID }()).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, repository.ErrNotUniqueName
		}

		return nil, fmt.Errorf("%w: %s", repository.ErrUpdateVocabulary, err.Error())
	}

	return v.ent2model(updated), err
}

func (v *Vocabulary) Delete(ctx context.Context, filter *repository.VocabularyFilter) error {
	_, err := v.client.Delete().Where(
		v.buildQuery(filter)...,
	).Exec(ctx)
	if err != nil {
		return errors.Join(repository.ErrDeleteVocabulary, err)
	}
	return nil
}

func (v *Vocabulary) Get(ctx context.Context, filter *repository.VocabularyFilter) ([]*model.Vocabulary, error) {
	entvoc, err := v.client.Query().Where(
		v.buildQuery(filter)...,
	).All(ctx)

	if err != nil {
		return nil, errors.Join(repository.ErrFindVocabulary, err)
	}

	vocabularies := make([]*model.Vocabulary, 0, len(entvoc))

	for _, voc := range entvoc {
		vocabularies = append(vocabularies, v.ent2model(voc))
	}

	return vocabularies, nil
}

func (v *Vocabulary) buildQuery(filter *repository.VocabularyFilter) []predicate.Vocabulary {
	var predicates = make([]predicate.Vocabulary, 0)
	// Filter by id
	if len(filter.ID) > 0 {
		predicates = append(predicates, vocabulary.IDIn(filter.ID...))
	}
	// Filter by parent id
	if len(filter.ParentID) > 0 {
		predicates = append(predicates, vocabulary.ParentIDIn(filter.ParentID...))
	}
	// Filter by name
	if len(filter.ParentID) > 0 {
		predicates = append(predicates, vocabulary.NameIn(filter.Name...))
	}

	return predicates
}

func (v *Vocabulary) ent2model(vocabulary *ent.Vocabulary) *model.Vocabulary {
	return &model.Vocabulary{
		ID: vocabulary.ID,
		Data: model.VocabularyData{
			Name:        vocabulary.Name,
			Title:       vocabulary.Title,
			Description: &vocabulary.Description,
			ParentID:    vocabulary.ParentID,
		},
	}
}
