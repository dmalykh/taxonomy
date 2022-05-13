package repository

import (
	"context"
	"tagservice/server/model"
)

type Relation interface {
	Create(ctx context.Context, relation ...*model.Relation) error
	Delete(ctx context.Context, filter *model.Relation) error
	// Get returns relation for specified arguments. Every relation should conform any of namespaceIds and any of entityIds and any of tagIds.
	// Not specified arguments ignored.
	Get(ctx context.Context, tagIds []uint64, namespaceIds []uint64, entityIds []uint64) ([]model.Relation, error)
}
