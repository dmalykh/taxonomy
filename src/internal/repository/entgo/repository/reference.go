package repository

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"errors"
	"fmt"
	"github.com/dmalykh/taxonomy/internal/repository/entgo/ent"
	"github.com/dmalykh/taxonomy/internal/repository/entgo/ent/predicate"
	"github.com/dmalykh/taxonomy/internal/repository/entgo/ent/reference"
	"github.com/dmalykh/taxonomy/taxonomy/model"
	"github.com/dmalykh/taxonomy/taxonomy/repository"
	"github.com/samber/lo"
)

type Reference struct {
	client *ent.ReferenceClient
}

func NewReference(client *ent.ReferenceClient) repository.Reference {
	return &Reference{
		client: client,
	}
}

func (r *Reference) Set(ctx context.Context, reference ...*repository.ReferenceModel) error {
	err := r.client.CreateBulk(func() []*ent.ReferenceCreate {
		create := make([]*ent.ReferenceCreate, 0, len(reference))

		for _, rel := range reference {
			create = append(create, r.client.Create().
				SetTermID(rel.TermID).
				SetNamespaceID(rel.NamespaceID).
				SetEntityID(string(rel.EntityID)),
			)
		}

		return create
	}()...).
		OnConflict(
			sql.ConflictColumns(`term_id`, `namespace_id`, `entity_id`),
			sql.ResolveWithNewValues(),
		).Exec(ctx)
	if err != nil {
		return fmt.Errorf(`%w: %s`, repository.ErrCreateReference, err.Error())
	}

	return nil
}

func (r *Reference) buildQuery(filter *repository.ReferenceFilter) []predicate.Reference {
	var predicates = make([]predicate.Reference, 0)

	// Filter by namespace
	if len(filter.NamespaceID) > 0 {
		predicates = append(predicates, reference.NamespaceIDIn(filter.NamespaceID...))
	}

	// Filter by terms
	if len(filter.TermID) > 0 {

		var groups = make([]predicate.Reference, 0, len(filter.TermID))
		for _, termID := range filter.TermID {
			terms := lo.Map[uint64, any](termID, func(item uint64, index int) any {
				return any(item)
			})
			groups = append(groups, func(s *sql.Selector) {
				s.Where(sql.In(
					s.C(reference.FieldEntityID),
					sql.Select(s.C(reference.FieldEntityID)).
						From(sql.Table(reference.Table)).
						Where(sql.In(
							s.C(reference.FieldTermID),
							terms...),
						),
				))
			})
		}
		predicates = append(predicates,
			reference.TermIDIn(lo.Uniq[uint64](lo.Flatten[uint64](filter.TermID))...),
			reference.And(groups...),
		)
	}

	// Filter by entity
	if len(filter.EntityID) > 0 {
		predicates = append(predicates, reference.EntityIDIn(func() []string {
			var ids = make([]string, 0, len(filter.EntityID))
			for _, id := range filter.EntityID {
				ids = append(ids, string(id))
			}
			return ids
		}()...))
	}

	// After id condition
	if filter.AfterID != nil {
		predicates = append(predicates, reference.IDGT(*filter.AfterID))
	}

	return predicates
}

func (r *Reference) Delete(ctx context.Context, filter *repository.ReferenceFilter) error {
	if len(filter.NamespaceID) == 0 {
		return repository.ErrWithoutNamespace
	}

	_, err := r.client.Delete().Where(
		r.buildQuery(filter)...,
	).Exec(ctx)
	if err != nil {
		return errors.Join(repository.ErrDeleteReferences, err)
	}

	return nil
}

func (r *Reference) Get(ctx context.Context, filter *repository.ReferenceFilter) ([]*repository.ReferenceModel, error) {
	if len(filter.NamespaceID) == 0 {
		return nil, repository.ErrWithoutNamespace
	}

	entreferences, err := r.client.Query().Where(
		reference.And(r.buildQuery(filter)...),
	).All(ctx)
	if err != nil {
		return nil, errors.Join(repository.ErrGetReference, err)
	}

	references := make([]*repository.ReferenceModel, 0, len(entreferences))
	for _, rel := range entreferences {
		references = append(references, &repository.ReferenceModel{
			ID:          rel.ID,
			TermID:      rel.TermID,
			NamespaceID: rel.NamespaceID,
			EntityID:    model.EntityID(rel.EntityID),
		})
	}

	return references, nil
}
