//nolint:nilnil
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

type Vocabulary struct {
	termService       taxonomy.Term
	vocabularyService taxonomy.Vocabulary
}

func (c *Vocabulary) Parent(ctx context.Context, obj *apimodel.Vocabulary) (*apimodel.Vocabulary, error) {
	if obj.ParentID == nil {
		return nil, nil
	}

	vocabulary, err := c.vocabularyService.GetByID(ctx, *(*uint)(unsafe.Pointer(obj.ParentID)))
	if err != nil {
		return nil, fmt.Errorf(`error %w to get vocabulary %d: %s`, taxonomy.ErrVocabularyNotFound, *obj.ParentID, err.Error())
	}

	gen := vocabulary2gen(vocabulary)

	return &gen, nil
}

func (c *Vocabulary) Children(ctx context.Context, obj *apimodel.Vocabulary) ([]*apimodel.Vocabulary, error) {
	if obj.ParentID == nil {
		return nil, nil
	}

	vocabularys, err := c.vocabularyService.Get(ctx, &model2.VocabularyFilter{
		ParentID: (*uint)(unsafe.Pointer(&obj.ParentID)),
	})
	if err != nil {
		return nil, gqlerror.Errorf(`error to get categories by entities %s`, err.Error())
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

func (c *Vocabulary) Terms(ctx context.Context, obj *apimodel.Vocabulary, first int64, after *string) (*genmodel.TermsConnection, error) {
	var afterID uint

	if after != nil {
		if err := cursor.Unmarshal(*after, &afterID); err != nil {
			return nil, fmt.Errorf(`error to unmarshal %q: %w`, *after, err)
		}
	}

	terms, err := c.termService.Get(ctx, &model2.TermFilter{
		VocabularyID: []uint{uint(obj.ID)},
		Limit:        uint(first + 1), // dirty hack to obtain HasNextPage
		AfterID:      &afterID,
	})
	if err != nil {
		return nil, fmt.Errorf(`error to get list %w`, err)
	}

	return termsConnection(terms, int(first)), nil
}
