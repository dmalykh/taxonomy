package repository

import (
	"context"
	"tagservice/server/model"
)

type Namespace interface {
	Create(ctx context.Context, name string) (model.Namespace, error)
	Update(ctx context.Context, id uint, name string) (model.Namespace, error)
	GetById(ctx context.Context, id uint) (model.Namespace, error)
	GetByName(ctx context.Context, name string) (model.Namespace, error)
	DeleteById(ctx context.Context, id uint) error
	GetList(ctx context.Context, limit, offset uint) ([]model.Namespace, error)
}
