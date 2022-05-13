package server

import (
	"context"
	"tagservice/server/model"
)

type Namespace interface {
	Create(ctx context.Context, name string) (model.Namespace, error)
	Update(ctx context.Context, uid uint64, name string) (model.Namespace, error)
	Delete(ctx context.Context, uid uint64) error
	GetList(ctx context.Context, limit, offset uint64) ([]model.Namespace, error)
	GetByName(ctx context.Context, namespace string) (model.Namespace, error)
}
