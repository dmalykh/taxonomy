package repository

import (
	"context"
	"tagservice/server/model"
)

type Category interface {
	Create(ctx context.Context, data *model.CategoryData) (model.Category, error)
	Update(ctx context.Context, id uint, data *model.CategoryData) (model.Category, error)
	DeleteById(ctx context.Context, id uint) error
	GetById(ctx context.Context, id uint) (model.Category, error)
	GetList(ctx context.Context, limit, offset uint) ([]model.Category, error)
}
