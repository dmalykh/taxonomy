package service

import (
	"context"
	"fmt"
	"github.com/dmalykh/taxonomy/taxonomy"

	apimodel "github.com/dmalykh/taxonomy/api/graphql/model"
)

type Entity struct {
	termService       taxonomy.Term
	vocabularyService taxonomy.Vocabulary
}

func (e *Entity) FindVocabularyByID(ctx context.Context, id int64) (apimodel.Vocabulary, error) {
	vocabulary, err := e.vocabularyService.GetByID(ctx, uint(id))
	if err != nil {
		return apimodel.Vocabulary{}, fmt.Errorf(`error %w to get vocabulary %d: %s`, taxonomy.ErrVocabularyNotFound, id, err.Error())
	}

	return vocabulary2gen(vocabulary), nil
}

func (e *Entity) FindTermByID(ctx context.Context, id int64) (apimodel.Term, error) {
	term, err := e.termService.GetByID(ctx, uint(id))
	if err != nil {
		return apimodel.Term{}, fmt.Errorf(`error %w to get term %d: %s`, taxonomy.ErrTermNotFound, id, err.Error())
	}

	return term2gen(term), nil
}
