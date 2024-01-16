package service

import (
	"context"
	"fmt"
	"github.com/dmalykh/taxonomy/taxonomy"
	model2 "github.com/dmalykh/taxonomy/taxonomy/model"
	"unsafe"

	"github.com/dmalykh/taxonomy/api/graphql/generated/genmodel"
	apimodel "github.com/dmalykh/taxonomy/api/graphql/model"
	"github.com/dmalykh/taxonomy/api/graphql/service/cursor"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type Query struct {
	termService       taxonomy.Term
	vocabularyService taxonomy.Vocabulary
}

func (q *Query) Term(ctx context.Context, id int64) (apimodel.Term, error) {
	term, err := q.termService.GetByID(ctx, uint(id))
	if err != nil {
		return apimodel.Term{}, fmt.Errorf(`error %w to get term %d: %s`, taxonomy.ErrTermNotFound, id, err.Error())
	}

	return term2gen(term), nil
}

func (q *Query) Terms(ctx context.Context, vocabularyID int64, name *string, first int64, after *string) (*genmodel.TermsConnection, error) {
	var afterID uint
	if after != nil {
		if err := cursor.Unmarshal(*after, &afterID); err != nil {
			return nil, fmt.Errorf(`error to unmarshal %q: %w`, *after, err)
		}
	}

	terms, err := q.termService.Get(ctx, &model2.TermFilter{
		VocabularyID: []uint{uint(vocabularyID)},
		Limit:        uint(first + 1), // dirty hack to obtain HasNextPage
		AfterID:      &afterID,
		Name:         name,
	})
	if err != nil {
		return nil, fmt.Errorf(`error to get list %w`, err)
	}

	return termsConnection(terms, int(first)), nil
}

func (q *Query) TermsByEntities(ctx context.Context, namespace string, entityID []int64) ([]*apimodel.Term, error) {
	terms, err := q.termService.GetTermsByEntities(ctx, namespace, int64stoUints(entityID)...)
	if err != nil {
		return nil, gqlerror.Errorf(`error to get terms by entities %s`, err.Error())
	}

	return func(terms []model2.Term) []*apimodel.Term {
		apiterms := make([]*apimodel.Term, len(terms))

		for i, term := range terms {
			term := term2gen(term)
			apiterms[i] = &term
		}

		return apiterms
	}(terms), nil
}

func (q *Query) Vocabulary(ctx context.Context, id int64) (apimodel.Vocabulary, error) {
	vocabulary, err := q.vocabularyService.GetByID(ctx, uint(id))

	return vocabulary2gen(vocabulary), err
}

func (q *Query) Categories(ctx context.Context, parentID *int64, name *string) ([]*apimodel.Vocabulary, error) {
	vocabularys, err := q.vocabularyService.Get(ctx, &model2.VocabularyFilter{
		ParentID: (*uint)(unsafe.Pointer(parentID)),
		Name:     name,
	})
	if err != nil {
		return nil, gqlerror.Errorf(`error to get categories by filter %s`, err.Error())
	}

	return func(vocabularys []model2.Vocabulary) []*apimodel.Vocabulary {
		apivocabularys := make([]*apimodel.Vocabulary, len(vocabularys))

		for i, vocabulary := range vocabularys {
			vocabulary := vocabulary2gen(vocabulary)
			apivocabularys[i] = &vocabulary
		}

		return apivocabularys
	}(vocabularys), nil
}
