package repository

import (
	"context"
	"tagservice/server/model"
)

type Namespace interface {
	Create(ctx context.Context, name string) (model.Namespace, error)
	Update(ctx context.Context, id uint64, name string) (model.Namespace, error)
	GetById(ctx context.Context, id uint64) (model.Namespace, error)
	GetByName(ctx context.Context, name string) (model.Namespace, error)
	DeleteById(ctx context.Context, id uint64) error
	GetList(ctx context.Context, limit, offset uint64) ([]model.Namespace, error)
}
