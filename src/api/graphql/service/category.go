package service

import (
	"context"
	"github.com/AlekSi/pointer"
	"github.com/dmalykh/tagservice/api/graphql/generated/genmodel"
	apimodel "github.com/dmalykh/tagservice/api/graphql/model"
	"github.com/dmalykh/tagservice/api/graphql/service/cursor"
	"github.com/dmalykh/tagservice/tagservice"
	"github.com/dmalykh/tagservice/tagservice/model"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"unsafe"
)

type Category struct {
	tagService      tagservice.Tag
	categoryService tagservice.Category
}

func (c *Category) Parent(ctx context.Context, obj *apimodel.Category) (*apimodel.Category, error) {
	if obj.ParentId == nil {
		return nil, nil
	}
	category, err := c.categoryService.GetById(ctx, *(*uint)(unsafe.Pointer(obj.ParentId)))
	var gen = category2gen(category)
	return &gen, err
}

func (c *Category) Children(ctx context.Context, obj *apimodel.Category) ([]*apimodel.Category, error) {
	if obj.ParentId == nil {
		return nil, nil
	}
	categorys, err := c.categoryService.GetList(ctx, &model.CategoryFilter{
		ParentId: (*uint)(unsafe.Pointer(&obj.ParentId)),
	})
	if err != nil {
		return nil, gqlerror.Errorf(`error to get categories by entities %w`, err)
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

func (c *Category) Tags(ctx context.Context, obj *apimodel.Category, first int64, after *string) (*genmodel.TagsConnection, error) {
	var afterId uint
	if after != nil {
		if err := cursor.Unmarshal(*after, &afterId); err != nil {
			return nil, err
		}
	}

	tags, err := c.tagService.GetList(ctx, &model.TagFilter{
		CategoryId: []uint{uint(obj.ID)},
		Limit:      pointer.ToUint(uint(first + 1)), // dirty hack to obtain HasNextPage
		AfterId:    &afterId,
	})
	if err != nil {
		return nil, err
	}
	return tagsConnection(tags, int(first)), nil
}
