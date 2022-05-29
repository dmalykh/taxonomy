package server

import (
	"context"
	"tagservice/server/model"
)

type Namespace interface {
	Create(ctx context.Context, name string) (model.Namespace, error)
	Update(ctx context.Context, uid uint, name string) (model.Namespace, error)
	Delete(ctx context.Context, uid uint) error
	GetList(ctx context.Context, limit, offset uint) ([]model.Namespace, error)
	GetByName(ctx context.Context, namespace string) (model.Namespace, error)
}
