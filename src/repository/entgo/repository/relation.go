package repository

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"fmt"
	"github.com/dmalykh/tagservice/repository/entgo/ent"
	"github.com/dmalykh/tagservice/repository/entgo/ent/relation"
	"github.com/dmalykh/tagservice/tagservice/model"
	"github.com/dmalykh/tagservice/tagservice/repository"
)

type Relation struct {
	client *ent.RelationClient
}

func NewRelation(client *ent.RelationClient) repository.Relation {
	return &Relation{
		client: client,
	}
}

func (r *Relation) Create(ctx context.Context, relation ...*model.Relation) error {
	added, err := r.client.CreateBulk(func() []*ent.RelationCreate {
		var create = make([]*ent.RelationCreate, 0, len(relation))
		for _, rel := range relation {
			create = append(create, r.client.Create().
				SetTagID(int(rel.TagId)).
				SetNamespaceID(int(rel.NamespaceId)).
				SetEntityID(int(rel.EntityId)))
		}
		return create
	}()...).Save(ctx)
	if err != nil {
		return fmt.Errorf(`%w: %s`, repository.ErrCreateRelation, err.Error())
	}
	if len(added) != len(relation) {
		return fmt.Errorf(`%w: internal error`, repository.ErrCreateRelation)
	}
	return nil
}

func (r *Relation) Delete(ctx context.Context, tagIds []uint, namespaceIds []uint, entityIds []uint) error {
	if len(entityIds) > 0 && len(namespaceIds) == 0 {
		return repository.ErrEntityWithoutNamespace
	}
	if _, err := r.client.Delete().Where(
		relation.And(
			// Remove by tags
			func(s *sql.Selector) {
				if len(tagIds) > 0 {
					s.Where(sql.In(relation.FieldTagID, func() []interface{} {
						var arr = make([]interface{}, len(tagIds))
						for i, id := range tagIds {
							arr[i] = id
						}
						return arr
					}()...))
				}
			},
			// Remove by namespaces
			func(s *sql.Selector) {
				if len(namespaceIds) > 0 {
					s.Where(sql.In(relation.FieldNamespaceID, func() []interface{} {
						var arr = make([]interface{}, len(namespaceIds))
						for i, id := range namespaceIds {
							arr[i] = id
						}
						return arr
					}()...))
				}
			},
			// Remove by entity
			func(s *sql.Selector) {
				if len(entityIds) > 0 && len(namespaceIds) > 0 {
					s.Where(sql.In(relation.FieldEntityID, func() []interface{} {
						var arr = make([]interface{}, len(entityIds))
						for i, id := range entityIds {
							arr[i] = id
						}
						return arr
					}()...))
				}
			},
		)).Exec(ctx); err != nil {
		return fmt.Errorf("%w (%+v, %+vx, %+v): %s", repository.ErrDeleteRelations, tagIds, namespaceIds, err, err.Error())
	}
	return nil
}

func (r *Relation) Get(ctx context.Context, tagIds []uint, namespaceIds []uint, entityIds []uint) ([]model.Relation, error) {
	if len(entityIds) > 0 && len(namespaceIds) == 0 {
		return nil, repository.ErrEntityWithoutNamespace
	}
	entrelations, err := r.client.Query().Where(
		relation.And(
			// By tags
			func(s *sql.Selector) {
				if len(tagIds) > 0 {
					s.Where(sql.In(relation.FieldTagID, func() []interface{} {
						var arr = make([]interface{}, len(tagIds))
						for i, id := range tagIds {
							arr[i] = id
						}
						return arr
					}()...))
				}
			},
			// By namespaces
			func(s *sql.Selector) {
				if len(namespaceIds) > 0 {
					s.Where(sql.In(relation.FieldNamespaceID, func() []interface{} {
						var arr = make([]interface{}, len(namespaceIds))
						for i, id := range namespaceIds {
							arr[i] = id
						}
						return arr
					}()...))
				}
			},
			// By entity
			func(s *sql.Selector) {
				if len(entityIds) > 0 && len(namespaceIds) > 0 {
					s.Where(sql.In(relation.FieldEntityID, func() []interface{} {
						var arr = make([]interface{}, len(entityIds))
						for i, id := range entityIds {
							arr[i] = id
						}
						return arr
					}()...))
				}
			},
		)).All(ctx)
	if err != nil {
		return nil, err
	}
	var relations = make([]model.Relation, 0, len(entrelations))
	for _, rel := range entrelations {
		relations = append(relations, model.Relation{
			TagId:       uint(rel.TagID),
			NamespaceId: uint(rel.NamespaceID),
			EntityId:    uint(rel.EntityID),
		})
	}
	return relations, nil
}
