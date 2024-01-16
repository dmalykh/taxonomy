package service

import (
	"context"
	"github.com/dmalykh/taxonomy/taxonomy"
	"github.com/dmalykh/taxonomy/taxonomy/model"

	"github.com/AlekSi/pointer"
	"github.com/dmalykh/taxonomy/api/graphql/generated/genmodel"
	apimodel "github.com/dmalykh/taxonomy/api/graphql/model"
	"github.com/dmalykh/taxonomy/api/graphql/service/cursor"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type Term struct {
	termService       taxonomy.Term
	vocabularyService taxonomy.Vocabulary
	namespaceService  taxonomy.Namespace
}

func (t *Term) Vocabulary(ctx context.Context, obj *apimodel.Term) (apimodel.Vocabulary, error) {
	vocabulary, err := t.vocabularyService.GetByID(ctx, uint(obj.VocabularyID))

	return vocabulary2gen(vocabulary), err
}

func (t *Term) Entities(ctx context.Context, obj *apimodel.Term, first int64, after *string, namespace []*string) (*genmodel.EntitiesConnection, error) { //nolint:lll
	var afterID uint
	if after != nil {
		if err := cursor.Unmarshal(*after, &afterID); err != nil {
			return nil, gqlerror.Errorf(`error to unmarshal cursor %s`, err.Error())
		}
	}

	// Get references
	references, err := t.termService.GetReferences(ctx, &model.EntityFilter{
		TermID: [][]uint{{uint(obj.ID)}},
		Namespace: func() []string {
			namespaces := make([]string, len(namespace))
			for i, ns := range namespace {
				namespaces[i] = *ns
			}

			return namespaces
		}(),
		AfterID: &afterID,
		Limit:   pointer.ToUint(uint(first + 1)),
	})
	if err != nil {
		return nil, gqlerror.Errorf(`error to get references %s`, err.Error())
	}

	// Get map of namespaces and their ids
	namespaces := func() map[uint]string {
		namespaces := make(map[uint]string, len(namespace))

		for _, ns := range namespace {
			namespaceModel, _ := t.namespaceService.GetByName(ctx, *ns)
			namespaces[namespaceModel.ID] = namespaceModel.Name
		}

		return namespaces
	}()

	return func(references []model.Reference, limit int) *genmodel.EntitiesConnection {
		var connection genmodel.EntitiesConnection

		if len(references) == 0 {
			return nil
		}

		if len(references) > limit {
			references = references[:limit]
		}

		connection.Edges = make([]genmodel.EntitiesEdge, len(references))
		for i, reference := range references {
			connection.Edges[i] = genmodel.EntitiesEdge{
				Cursor: cursor.Marshal(reference.ID),
				Node: &genmodel.EntityNode{
					Namespace: namespaces[reference.NamespaceID],
					ID:        int64(reference.EntityID),
				},
			}
		}

		connection.PageInfo = genmodel.PageInfo{
			StartCursor: connection.Edges[0].Cursor,
			EndCursor:   connection.Edges[len(connection.Edges)-1].Cursor,
			HasNextPage: pointer.ToBool(len(references) > limit),
		}

		return &connection
	}(references, int(first)), nil
}
