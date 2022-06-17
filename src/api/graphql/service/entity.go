package service

import (
	"context"
	"fmt"

	apimodel "github.com/dmalykh/tagservice/api/graphql/model"
	"github.com/dmalykh/tagservice/tagservice"
)

//goland:noinspection GoUnnecessarilyExportedIdentifiers
type Entity struct {
	tagService      tagservice.Tag
	categoryService tagservice.Category
}

func (e *Entity) FindCategoryByID(ctx context.Context, id int64) (apimodel.Category, error) {
	category, err := e.categoryService.GetByID(ctx, uint(id))
	if err != nil {
		return apimodel.Category{}, fmt.Errorf(`error %w to get category %d: %s`, tagservice.ErrCategoryNotFound, id, err.Error())
	}

	return category2gen(category), nil
}

func (e *Entity) FindTagByID(ctx context.Context, id int64) (apimodel.Tag, error) {
	tag, err := e.tagService.GetByID(ctx, uint(id))
	if err != nil {
		return apimodel.Tag{}, fmt.Errorf(`error %w to get tag %d: %s`, tagservice.ErrTagNotFound, id, err.Error())
	}

	return tag2gen(tag), nil
}
