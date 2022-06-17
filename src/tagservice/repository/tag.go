package repository

import (
	"context"
	"errors"

	"github.com/dmalykh/tagservice/tagservice/model"
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
	DeleteByID(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (model.Tag, error)
	GetByName(ctx context.Context, name string) ([]model.Tag, error)
	GetList(ctx context.Context, filter *model.TagFilter) ([]model.Tag, error)
}
