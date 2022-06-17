package service

import (
	"context"
	"fmt"
	"unsafe"

	"github.com/dmalykh/tagservice/api/graphql/generated/genmodel"
	apimodel "github.com/dmalykh/tagservice/api/graphql/model"
	"github.com/dmalykh/tagservice/api/graphql/service/cursor"
	"github.com/dmalykh/tagservice/tagservice"
	"github.com/dmalykh/tagservice/tagservice/model"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

//goland:noinspection GoUnnecessarilyExportedIdentifiers
type Query struct {
	tagService      tagservice.Tag
	categoryService tagservice.Category
}

func (q *Query) Tag(ctx context.Context, id int64) (apimodel.Tag, error) {
	tag, err := q.tagService.GetByID(ctx, uint(id))
	if err != nil {
		return apimodel.Tag{}, fmt.Errorf(`error %w to get tag %d: %s`, tagservice.ErrTagNotFound, id, err.Error())
	}

	return tag2gen(tag), nil
}

func (q *Query) Tags(ctx context.Context, categoryID int64, name *string, first int64, after *string) (*genmodel.TagsConnection, error) {
	var afterID uint
	if after != nil {
		if err := cursor.Unmarshal(*after, &afterID); err != nil {
			return nil, fmt.Errorf(`error to unmarshal %q: %w`, *after, err)
		}
	}

	tags, err := q.tagService.GetList(ctx, &model.TagFilter{
		CategoryID: []uint{uint(categoryID)},
		Limit:      uint(first + 1), // dirty hack to obtain HasNextPage
		AfterID:    &afterID,
		Name:       name,
	})
	if err != nil {
		return nil, fmt.Errorf(`error to get list %w`, err)
	}

	return tagsConnection(tags, int(first)), nil
}

func (q *Query) TagsByEntities(ctx context.Context, namespace string, entityID []int64) ([]*apimodel.Tag, error) {
	tags, err := q.tagService.GetTagsByEntities(ctx, namespace, int64stoUints(entityID)...)
	if err != nil {
		return nil, gqlerror.Errorf(`error to get tags by entities %s`, err.Error())
	}

	return func(tags []model.Tag) []*apimodel.Tag {
		apitags := make([]*apimodel.Tag, len(tags))

		for i, tag := range tags {
			tag := tag2gen(tag)
			apitags[i] = &tag
		}

		return apitags
	}(tags), nil
}

func (q *Query) Category(ctx context.Context, id int64) (apimodel.Category, error) {
	category, err := q.categoryService.GetByID(ctx, uint(id))

	return category2gen(category), err
}

func (q *Query) Categories(ctx context.Context, parentID *int64, name *string) ([]*apimodel.Category, error) {
	categorys, err := q.categoryService.GetList(ctx, &model.CategoryFilter{
		ParentID: (*uint)(unsafe.Pointer(parentID)),
		Name:     name,
	})
	if err != nil {
		return nil, gqlerror.Errorf(`error to get categories by filter %s`, err.Error())
	}

	return func(categorys []model.Category) []*apimodel.Category {
		apicategorys := make([]*apimodel.Category, len(categorys))

		for i, category := range categorys {
			category := category2gen(category)
			apicategorys[i] = &category
		}

		return apicategorys
	}(categorys), nil
}
