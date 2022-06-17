//nolint:nilnil
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
type Category struct {
	tagService      tagservice.Tag
	categoryService tagservice.Category
}

func (c *Category) Parent(ctx context.Context, obj *apimodel.Category) (*apimodel.Category, error) {
	if obj.ParentID == nil {
		return nil, nil
	}

	category, err := c.categoryService.GetByID(ctx, *(*uint)(unsafe.Pointer(obj.ParentID)))
	if err != nil {
		return nil, fmt.Errorf(`error %w to get category %d: %s`, tagservice.ErrCategoryNotFound, *obj.ParentID, err.Error())
	}

	gen := category2gen(category)

	return &gen, nil
}

func (c *Category) Children(ctx context.Context, obj *apimodel.Category) ([]*apimodel.Category, error) {
	if obj.ParentID == nil {
		return nil, nil
	}

	categorys, err := c.categoryService.GetList(ctx, &model.CategoryFilter{
		ParentID: (*uint)(unsafe.Pointer(&obj.ParentID)),
	})
	if err != nil {
		return nil, gqlerror.Errorf(`error to get categories by entities %s`, err.Error())
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

func (c *Category) Tags(ctx context.Context, obj *apimodel.Category, first int64, after *string) (*genmodel.TagsConnection, error) {
	var afterID uint

	if after != nil {
		if err := cursor.Unmarshal(*after, &afterID); err != nil {
			return nil, fmt.Errorf(`error to unmarshal %q: %w`, *after, err)
		}
	}

	tags, err := c.tagService.GetList(ctx, &model.TagFilter{
		CategoryID: []uint{uint(obj.ID)},
		Limit:      uint(first + 1), // dirty hack to obtain HasNextPage
		AfterID:    &afterID,
	})
	if err != nil {
		return nil, fmt.Errorf(`error to get list %w`, err)
	}

	return tagsConnection(tags, int(first)), nil
}
