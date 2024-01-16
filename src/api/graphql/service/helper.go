package service

import (
	model2 "github.com/dmalykh/taxonomy/taxonomy/model"
	"unsafe"

	"github.com/AlekSi/pointer"
	"github.com/dmalykh/taxonomy/api/graphql/generated/genmodel"
	apimodel "github.com/dmalykh/taxonomy/api/graphql/model"
	"github.com/dmalykh/taxonomy/api/graphql/service/cursor"
)

func term2gen(term model2.Term) apimodel.Term {
	return apimodel.Term{
		ID:           int64(term.ID),
		Name:         term.Data.Name,
		Title:        &term.Data.Title,
		Description:  &term.Data.Description,
		VocabularyID: int64(term.Data.VocabularyID),
	}
}

func vocabulary2gen(vocabulary model2.Vocabulary) apimodel.Vocabulary {
	return apimodel.Vocabulary{
		ID:          int64(vocabulary.ID),
		Name:        vocabulary.Data.Name,
		Title:       vocabulary.Data.Title,
		Description: vocabulary.Data.Description,
		ParentID:    (*int64)(unsafe.Pointer(vocabulary.Data.ParentID)),
	}
}

func int64stoUints(ints []int64) []uint {
	uints := make([]uint, 0, len(ints))
	for _, i := range ints {
		uints = append(uints, uint(i))
	}

	return uints
}

func termsConnection(terms []model2.Term, limit int) *genmodel.TermsConnection {
	var connection genmodel.TermsConnection

	if len(terms) == 0 {
		return nil
	}

	if len(terms) > limit {
		terms = terms[:limit]
	}

	connection.Edges = make([]genmodel.TermsEdge, len(terms))

	for i, term := range terms {
		gen := term2gen(term)

		connection.Edges[i] = genmodel.TermsEdge{
			Cursor: cursor.Marshal(term.ID),
			Node:   &gen,
		}
	}

	connection.PageInfo = genmodel.PageInfo{
		StartCursor: connection.Edges[0].Cursor,
		EndCursor:   connection.Edges[len(connection.Edges)-1].Cursor,
		HasNextPage: pointer.ToBool(len(terms) > limit),
	}

	return &connection
}
