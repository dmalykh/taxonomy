package service

import (
	"context"
	"github.com/AlekSi/pointer"
	"github.com/dmalykh/tagservice/api/graphql/generated/genmodel"
	apimodel "github.com/dmalykh/tagservice/api/graphql/model"
	"github.com/dmalykh/tagservice/api/graphql/service/cursor"
	"github.com/dmalykh/tagservice/tagservice"
	"github.com/dmalykh/tagservice/tagservice/model"
)

type Tag struct {
	tagService       tagservice.Tag
	categoryService  tagservice.Category
	namespaceService tagservice.Namespace
}

func (t *Tag) Category(ctx context.Context, obj *apimodel.Tag) (apimodel.Category, error) {
	category, err := t.categoryService.GetById(ctx, uint(obj.CategoryId))
	return category2gen(category), err
}

func (t *Tag) Entities(ctx context.Context, obj *apimodel.Tag, first int64, after *string, namespace []*string) (*genmodel.EntitiesConnection, error) {
	var afterId uint
	if after != nil {
		if err := cursor.Unmarshal(*after, &afterId); err != nil {
			return nil, err
		}
	}

	// Get relations
	relations, err := t.tagService.GetRelations(ctx, &model.RelationFilter{
		TagId: []uint{uint(obj.ID)},
		Namespace: func() []string {
			var namespaces = make([]string, len(namespace))
			for i, ns := range namespace {
				namespaces[i] = *ns
			}
			return namespaces
		}(),
		AfterId: &afterId,
		Limit:   pointer.ToUint(uint(first + 1)),
	})
	if err != nil {
		return nil, err //@TODO
	}

	// Get map of namespaces and their ids
	var namespaces = func() map[uint]string {
		var namespaces = make(map[uint]string, len(namespace))
		for _, ns := range namespace {
			namespaceModel, _ := t.namespaceService.GetByName(ctx, *ns)
			namespaces[namespaceModel.Id] = namespaceModel.Name
		}
		return namespaces
	}()

	return func(relations []model.Relation, limit int) *genmodel.EntitiesConnection {
		var connection genmodel.EntitiesConnection
		if len(relations) == 0 {
			return nil
		}
		if len(relations) > limit {
			relations = relations[:limit]
		}
		connection.Edges = make([]genmodel.EntitiesEdge, len(relations))
		for i, relation := range relations {
			connection.Edges[i] = genmodel.EntitiesEdge{
				Cursor: cursor.Marshal(relation.Id),
				Node: &genmodel.EntityNode{
					Namespace: namespaces[relation.NamespaceId],
					ID:        int64(relation.EntityId),
				},
			}
		}
		connection.PageInfo = genmodel.PageInfo{
			StartCursor: connection.Edges[0].Cursor,
			EndCursor:   connection.Edges[len(connection.Edges)-1].Cursor,
			HasNextPage: pointer.ToBool(len(relations) > limit),
		}
		return &connection
	}(relations, int(first)), nil

}
