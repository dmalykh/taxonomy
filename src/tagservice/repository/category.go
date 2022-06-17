package repository

import (
	"context"
	"errors"

	"github.com/dmalykh/tagservice/tagservice/model"
)

var (
	ErrCreateCategory = errors.New(`failed to create category`)
	ErrUpdateCategory = errors.New(`failed to update category`)
	ErrNotUniqueName  = errors.New(`category's name and parent must be unique`)
	ErrFindCategory   = errors.New(`failed to find category`)
	ErrDeleteCategory = errors.New(`failed to delete category`)
)

type Category interface {
	Create(ctx context.Context, data *model.CategoryData) (model.Category, error)
	Update(ctx context.Context, id uint, data *model.CategoryData) (model.Category, error)
	DeleteByID(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (model.Category, error)
	GetList(ctx context.Context, filter *model.CategoryFilter) ([]model.Category, error)
}
