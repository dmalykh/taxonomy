package service

import (
	"context"
	"github.com/dmalykh/taxonomy/taxonomy"
	model2 "github.com/dmalykh/taxonomy/taxonomy/model"
	"unsafe"

	"github.com/AlekSi/pointer"
	"github.com/dmalykh/taxonomy/api/graphql/generated/genmodel"
	apimodel "github.com/dmalykh/taxonomy/api/graphql/model"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type Mutation struct {
	termService       taxonomy.Term
	vocabularyService taxonomy.Vocabulary
}

func (m *Mutation) CreateTerm(ctx context.Context, input genmodel.TermInput) (apimodel.Term, error) {
	term, err := m.termService.Create(ctx, &model2.TermData{
		Name:         input.Name,
		Title:        input.Title,
		VocabularyID: uint(input.VocabularyID),
		Description:  *(input.Description),
	})
	if err != nil {
		return apimodel.Term{}, gqlerror.Errorf(`error to create term %s`, err.Error())
	}

	return term2gen(term), nil
}

func (m *Mutation) UpdateTerm(ctx context.Context, id int64, input genmodel.TermInput) (apimodel.Term, error) {
	term, err := m.termService.Update(ctx, uint(id), &model2.TermData{
		Name:         input.Name,
		Title:        input.Title,
		VocabularyID: uint(input.VocabularyID),
		Description:  *(input.Description),
	})
	if err != nil {
		return apimodel.Term{}, gqlerror.Errorf(`error to update term %s`, err.Error())
	}

	return term2gen(term), nil
}

func (m *Mutation) Set(ctx context.Context, termID []int64, namespace string, entityID []int64) (*bool, error) {
	entitiesID := int64stoUints(entityID)
	for _, id := range termID {
		if err := m.termService.SetReference(ctx, uint(id), namespace, entitiesID...); err != nil {
			return nil, gqlerror.Errorf(`error to set reference %d %s %q term %s`, id, namespace, entitiesID, err.Error())
		}
	}

	return pointer.ToBool(true), nil
}

func (m *Mutation) Unset(ctx context.Context, termID []int64, namespace string, entityID []int64) (*bool, error) {
	entitiesID := int64stoUints(entityID)
	for _, id := range termID {
		if err := m.termService.SetReference(ctx, uint(id), namespace, entitiesID...); err != nil {
			return nil, gqlerror.Errorf(`error to unset reference %d %s %q term %s`, id, namespace, entitiesID, err.Error())
		}
	}

	return pointer.ToBool(true), nil
}

func (m *Mutation) CreateVocabulary(ctx context.Context, input genmodel.VocabularyInput) (apimodel.Vocabulary, error) {
	vocabulary, err := m.vocabularyService.Create(ctx, &model2.VocabularyData{
		Name:        input.Name,
		Title:       input.Title,
		ParentID:    (*uint)(unsafe.Pointer(input.ParentID)),
		Description: input.Description,
	})
	if err != nil {
		return apimodel.Vocabulary{}, gqlerror.Errorf(`error to create vocabulary %s`, err.Error())
	}

	return vocabulary2gen(vocabulary), nil
}

func (m *Mutation) UpdateVocabulary(ctx context.Context, id int64, input genmodel.VocabularyInput) (apimodel.Vocabulary, error) {
	vocabulary, err := m.vocabularyService.Update(ctx, uint(id), &model2.VocabularyData{
		Name:        input.Name,
		Title:       input.Title,
		ParentID:    (*uint)(unsafe.Pointer(input.ParentID)),
		Description: input.Description,
	})
	if err != nil {
		return apimodel.Vocabulary{}, gqlerror.Errorf(`error to update vocabulary %s`, err.Error())
	}

	return vocabulary2gen(vocabulary), nil
}
