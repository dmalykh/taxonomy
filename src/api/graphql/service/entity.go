package service

import (
	"context"
	apimodel "github.com/dmalykh/tagservice/api/graphql/model"
	"github.com/dmalykh/tagservice/tagservice"
)

type Entity struct {
	tagService      tagservice.Tag
	categoryService tagservice.Category
}

func (e *Entity) FindCategoryByID(ctx context.Context, id int64) (apimodel.Category, error) {
	category, err := e.categoryService.GetById(ctx, uint(id))
	return category2gen(category), err
}

func (e *Entity) FindTagByID(ctx context.Context, id int64) (apimodel.Tag, error) {
	tag, err := e.tagService.GetById(ctx, uint(id))
	return tag2gen(tag), err
}
