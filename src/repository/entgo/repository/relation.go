package repository

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/dmalykh/tagservice/repository/entgo/ent"
	"github.com/dmalykh/tagservice/repository/entgo/ent/predicate"
	"github.com/dmalykh/tagservice/repository/entgo/ent/relation"
	"github.com/dmalykh/tagservice/tagservice/model"
	"github.com/dmalykh/tagservice/tagservice/repository"
)

//goland:noinspection GoUnnecessarilyExportedIdentifiers
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
		create := make([]*ent.RelationCreate, 0, len(relation))

		for _, rel := range relation {
			create = append(create, r.client.Create().
				SetTagID(int(rel.TagID)).
				SetNamespaceID(int(rel.NamespaceID)).
				SetEntityID(int(rel.EntityID)))
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

func (r *Relation) Delete(ctx context.Context, tagIds []uint, namespaceIds []uint, entityIds []uint) error { //nolint:cyclop
	if len(entityIds) > 0 && len(namespaceIds) == 0 {
		return repository.ErrEntityWithoutNamespace
	}

	if _, err := r.client.Delete().Where(
		relation.And(
			// Remove by tags
			func(s *sql.Selector) {
				if len(tagIds) > 0 {
					s.Where(sql.In(relation.FieldTagID, func() []interface{} {
						arr := make([]interface{}, len(tagIds))
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
						arr := make([]interface{}, len(namespaceIds))
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
						arr := make([]interface{}, len(entityIds))
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

//goland:noinspection GoRedundantConversion
func (r *Relation) Get(ctx context.Context, filter *model.RelationFilter) ([]model.Relation, error) { //nolint:gocognit,cyclop
	if len(filter.EntityID) > 0 && len(filter.Namespace) == 0 {
		return nil, repository.ErrEntityWithoutNamespace
	}

	entrelations, err := r.client.Query().Where(
		relation.And(func() []predicate.Relation {
			predicates := make([]predicate.Relation, 0, 5)
			// By tags
			if len(filter.TagID) > 0 {
				predicates = append(predicates, predicate.Relation(func(s *sql.Selector) {
					s.Where(sql.InInts(relation.FieldTagID, func() []int {
						arr := make([]int, 0, len(filter.TagID))
						for _, group := range filter.TagID {
							for _, id := range group {
								arr = append(arr, int(id))
							}
						}

						return arr
					}()...))
				}))
			}
			// By namespaces
			if len(filter.Namespace) > 0 {
				//goland:noinspection GoRedundantConversion
				predicates = append(predicates, predicate.Relation(func(s *sql.Selector) {
					s.Where(sql.In(relation.FieldNamespaceID, func() []interface{} {
						arr := make([]interface{}, len(filter.Namespace))
						for i, id := range filter.Namespace {
							arr[i] = id
						}

						return arr
					}()...))
				}))
			}
			// By entity
			if len(filter.EntityID) > 0 && len(filter.Namespace) > 0 {
				predicates = append(predicates, predicate.Relation(func(s *sql.Selector) {
					s.Where(sql.InInts(relation.FieldEntityID, func() []int {
						arr := make([]int, len(filter.EntityID))
						for i, id := range filter.EntityID {
							arr[i] = int(id)
						}

						return arr
					}()...))
				}))
			}
			// After id condition
			if filter.AfterID != nil {
				predicates = append(predicates, relation.IDGT(int(*filter.AfterID)))
			}

			return predicates
		}()...),
	).All(ctx)
	if err != nil {
		return nil, fmt.Errorf(`%w: %s`, repository.ErrGetRelation, err.Error())
	}

	relations := make([]model.Relation, 0, len(entrelations))
	for _, rel := range entrelations {
		relations = append(relations, model.Relation{
			ID:          uint(rel.ID),
			TagID:       uint(rel.TagID),
			NamespaceID: uint(rel.NamespaceID),
			EntityID:    uint(rel.EntityID),
		})
	}

	return relations, nil
}
