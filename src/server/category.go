package server

import (
	"context"
	"tagservice/server/model"
)

type Category interface {
	Create(ctx context.Context, data *model.CategoryData) (model.Category, error)
	Update(ctx context.Context, id uint64, data *model.CategoryData) (model.Category, error)
	Delete(ctx context.Context, id uint64) error
	GetById(ctx context.Context, id uint64) (model.Category, error)
	GetList(ctx context.Context, limit, offset uint64) ([]model.Category, error)
}
