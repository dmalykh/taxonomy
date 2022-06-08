package repository

import (
	"context"
	"errors"
	"tagservice/server/model"
)

var (
	ErrCreateTag = errors.New(`failed to create tag`)
	ErrUpdateTag = errors.New(`failed to update tag`)
	ErrFindTag   = errors.New(`failed to find tag`)
	ErrDeleteTag = errors.New(`failed to delete tag`)
)

type Tag interface {
	Create(ctx context.Context, data *model.TagData) (model.Tag, error)
	Update(ctx context.Context, id uint, data *model.TagData) (model.Tag, error)
	DeleteById(ctx context.Context, id uint) error
	GetById(ctx context.Context, id uint) (model.Tag, error)
	GetByName(ctx context.Context, name string) ([]model.Tag, error)
	GetByFilter(ctx context.Context, filter model.TagFilter, limit, offset uint) ([]model.Tag, error)
}
