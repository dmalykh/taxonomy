package service

import (
	"context"
	"github.com/dmalykh/tagservice/api/graphql/generated/genmodel"
	apimodel "github.com/dmalykh/tagservice/api/graphql/model"
	"github.com/dmalykh/tagservice/api/graphql/service/cursor"
	"github.com/dmalykh/tagservice/tagservice"
	"github.com/dmalykh/tagservice/tagservice/model"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"unsafe"
)

type Query struct {
	tagService      tagservice.Tag
	categoryService tagservice.Category
}

func (q *Query) Tag(ctx context.Context, id int64) (apimodel.Tag, error) {
	tag, err := q.tagService.GetById(ctx, uint(id))
	return tag2gen(tag), err
}

func (q *Query) Tags(ctx context.Context, categoryID int64, name *string, first int64, after *string) (*genmodel.TagsConnection, error) {
	var afterId uint
	if after != nil {
		if err := cursor.Unmarshal(*after, &afterId); err != nil {
			return nil, err
		}
	}

	tags, err := q.tagService.GetList(ctx, &model.TagFilter{
		CategoryId: []uint{uint(categoryID)},
		Limit:      uint(first + 1), // dirty hack to obtain HasNextPage
		AfterId:    &afterId,
		Name:       name,
	})
	if err != nil {
		return nil, err
	}
	return tagsConnection(tags, int(first)), nil
}

func (q *Query) TagsByEntities(ctx context.Context, namespace string, entityID []int64) ([]*apimodel.Tag, error) {
	tags, err := q.tagService.GetTagsByEntities(ctx, namespace, int64stoUints(entityID)...)
	if err != nil {
		return nil, gqlerror.Errorf(`error to get tags by entities %s`, err.Error())
	}
	return func(tags []model.Tag) []*apimodel.Tag {
		var apitags = make([]*apimodel.Tag, len(tags))
		for i, tag := range tags {
			var tag = tag2gen(tag)
			apitags[i] = &tag
		}
		return apitags
	}(tags), nil
}

func (q *Query) Category(ctx context.Context, id int64) (apimodel.Category, error) {
	category, err := q.categoryService.GetById(ctx, uint(id))
	return category2gen(category), err
}

func (q *Query) Categories(ctx context.Context, parentID *int64, name *string) ([]*apimodel.Category, error) {
	categorys, err := q.categoryService.GetList(ctx, &model.CategoryFilter{
		ParentId: (*uint)(unsafe.Pointer(parentID)),
		Name:     name,
	})
	if err != nil {
		return nil, gqlerror.Errorf(`error to get categories by filter %s`, err.Error())
	}
	return func(categorys []model.Category) []*apimodel.Category {
		var apicategorys = make([]*apimodel.Category, len(categorys))
		for i, category := range categorys {
			var category = category2gen(category)
			apicategorys[i] = &category
		}
		return apicategorys
	}(categorys), nil
}
