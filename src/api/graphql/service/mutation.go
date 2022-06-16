package service

import (
	"context"
	"github.com/AlekSi/pointer"
	"github.com/dmalykh/tagservice/api/graphql/generated/genmodel"
	apimodel "github.com/dmalykh/tagservice/api/graphql/model"
	"github.com/dmalykh/tagservice/tagservice"
	"github.com/dmalykh/tagservice/tagservice/model"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"unsafe"
)

type Mutation struct {
	tagService      tagservice.Tag
	categoryService tagservice.Category
}

func (m *Mutation) CreateTag(ctx context.Context, input genmodel.TagInput) (apimodel.Tag, error) {
	tag, err := m.tagService.Create(ctx, &model.TagData{
		Name:        input.Name,
		Title:       input.Title,
		CategoryId:  uint(input.CategoryID),
		Description: *(input.Description),
	})
	if err != nil {
		return apimodel.Tag{}, gqlerror.Errorf(`error to create tag %s`, err.Error())
	}
	return tag2gen(tag), nil
}

func (m *Mutation) UpdateTag(ctx context.Context, id int64, input genmodel.TagInput) (apimodel.Tag, error) {
	tag, err := m.tagService.Update(ctx, uint(id), &model.TagData{
		Name:        input.Name,
		Title:       input.Title,
		CategoryId:  uint(input.CategoryID),
		Description: *(input.Description),
	})
	if err != nil {
		return apimodel.Tag{}, gqlerror.Errorf(`error to update tag %s`, err.Error())
	}
	return tag2gen(tag), nil
}

func (m *Mutation) Set(ctx context.Context, tagID []int64, namespace string, entityID []int64) (*bool, error) {
	var entitiesId = int64stoUints(entityID)
	for _, id := range tagID {
		if err := m.tagService.SetRelation(ctx, uint(id), namespace, entitiesId...); err != nil {
			return nil, gqlerror.Errorf(`error to set relation %d %s %q tag %s`, id, namespace, entitiesId, err.Error())
		}
	}
	return pointer.ToBool(true), nil
}

func (m *Mutation) Unset(ctx context.Context, tagID []int64, namespace string, entityID []int64) (*bool, error) {
	var entitiesId = int64stoUints(entityID)
	for _, id := range tagID {
		if err := m.tagService.SetRelation(ctx, uint(id), namespace, entitiesId...); err != nil {
			return nil, gqlerror.Errorf(`error to unset relation %d %s %q tag %s`, id, namespace, entitiesId, err.Error())
		}
	}
	return pointer.ToBool(true), nil
}

func (m *Mutation) CreateCategory(ctx context.Context, input genmodel.CategoryInput) (apimodel.Category, error) {
	category, err := m.categoryService.Create(ctx, &model.CategoryData{
		Name:        input.Name,
		Title:       input.Title,
		ParentId:    (*uint)(unsafe.Pointer(input.ParentID)),
		Description: input.Description,
	})
	if err != nil {
		return apimodel.Category{}, gqlerror.Errorf(`error to create category %s`, err.Error())
	}
	return category2gen(category), nil
}

func (m *Mutation) UpdateCategory(ctx context.Context, id int64, input genmodel.CategoryInput) (apimodel.Category, error) {
	category, err := m.categoryService.Update(ctx, uint(id), &model.CategoryData{
		Name:        input.Name,
		Title:       input.Title,
		ParentId:    (*uint)(unsafe.Pointer(input.ParentID)),
		Description: input.Description,
	})
	if err != nil {
		return apimodel.Category{}, gqlerror.Errorf(`error to update category %s`, err.Error())
	}
	return category2gen(category), nil
}
