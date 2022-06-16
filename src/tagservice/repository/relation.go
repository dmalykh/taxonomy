package repository

import (
	"context"
	"errors"
	"github.com/dmalykh/tagservice/tagservice/model"
)

var (
	ErrCreateRelation         = errors.New(`failed to create relation`)
	ErrEntityWithoutNamespace = errors.New(`namespace required when requiring entity`)
	ErrDeleteRelations        = errors.New(`failed to delete relation`)
)

type Relation interface {
	Create(ctx context.Context, relation ...*model.Relation) error
	Delete(ctx context.Context, tagIds []uint, namespaceIds []uint, entityIds []uint) error
	Get(ctx context.Context, filter *model.RelationFilter) ([]model.Relation, error)
}
