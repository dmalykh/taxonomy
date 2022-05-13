package repository

import (
	"context"
	"tagservice/server/model"
)

type Tag interface {
	Create(ctx context.Context, data *model.TagData) (model.Tag, error)
	Update(ctx context.Context, id uint64, data *model.TagData) (model.Tag, error)
	DeleteById(ctx context.Context, id uint64) error
	GetById(ctx context.Context, id uint64) (model.Tag, error)
	GetByName(ctx context.Context, name string) (model.Tag, error)
	GetByFilter(ctx context.Context, filter model.TagFilter, limit, offset uint64) ([]model.Tag, error)
}
