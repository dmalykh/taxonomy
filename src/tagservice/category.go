package tagservice

import (
	"context"
	"errors"
	"github.com/dmalykh/tagservice/tagservice/model"
)

var ErrCategoryNotFound = errors.New(`category not found`)
var ErrCategoryNotCreated = errors.New(`category had not created`)
var ErrCategoryHasTags = errors.New(`category has tags, but should be empty`)
var ErrCategoryNotUpdated = errors.New(`category had not updated`)

type Category interface {
	Create(ctx context.Context, data *model.CategoryData) (model.Category, error)
	Update(ctx context.Context, id uint, data *model.CategoryData) (model.Category, error)
	Delete(ctx context.Context, id uint) error
	GetById(ctx context.Context, id uint) (model.Category, error)
	GetList(ctx context.Context, limit, offset uint) ([]model.Category, error)
}
