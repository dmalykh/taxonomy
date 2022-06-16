package repository

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"fmt"
	"github.com/dmalykh/tagservice/repository/entgo/ent"
	"github.com/dmalykh/tagservice/repository/entgo/ent/category"
	"github.com/dmalykh/tagservice/tagservice/model"
	"github.com/dmalykh/tagservice/tagservice/repository"
	"unsafe"
)

type Category struct {
	client *ent.CategoryClient
}

func NewCategory(client *ent.CategoryClient) repository.Category {
	return &Category{
		client: client,
	}
}

func (c *Category) Create(ctx context.Context, data *model.CategoryData) (model.Category, error) {
	ns, err := c.client.Create().
		SetName(data.Name).
		SetTitle(data.Title).
		SetNillableDescription(data.Description).
		SetNillableParentID(func() *int { return (*int)(unsafe.Pointer(data.ParentId)) }()).
		Save(ctx)

	if err != nil {
		if ent.IsConstraintError(err) {
			return model.Category{}, repository.ErrNotUniqueName
		}
		return model.Category{}, fmt.Errorf("%w: %s", repository.ErrCreateCategory, err.Error())
	}
	return c.ent2model(ns), nil
}

func (c *Category) Update(ctx context.Context, id uint, data *model.CategoryData) (model.Category, error) {
	category, err := c.client.UpdateOneID(int(id)).
		SetName(data.Name).
		SetTitle(data.Title).
		SetNillableDescription(data.Description).
		ClearParentID().
		SetNillableParentID(func() *int { return (*int)(unsafe.Pointer(data.ParentId)) }()).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return model.Category{}, repository.ErrNotUniqueName
		}
		return model.Category{}, fmt.Errorf("%w: %s", repository.ErrUpdateCategory, err.Error())
	}
	return c.ent2model(category), err
}

func (c *Category) DeleteById(ctx context.Context, id uint) error {
	if err := c.client.DeleteOneID(int(id)).Exec(ctx); err != nil {
		return fmt.Errorf("%w (%d): %s", repository.ErrDeleteCategory, id, err.Error())
	}
	return nil
}

func (c *Category) GetById(ctx context.Context, id uint) (model.Category, error) {
	ns, err := c.client.Get(ctx, int(id))
	if err != nil {
		return model.Category{}, fmt.Errorf("%w (%d): %s", repository.ErrFindCategory, id, err.Error())
	}
	return c.ent2model(ns), err
}

func (c *Category) GetList(ctx context.Context, filter *model.CategoryFilter) ([]model.Category, error) {
	entcategories, err := c.client.Query().Where(func(s *sql.Selector) {
		// Filter by parent id
		if filter.ParentId != nil {
			s.Where(sql.EQ(category.FieldParentID, *filter.ParentId))
		}
		// Filter by name
		if filter.Name != nil {
			s.Where(sql.EQ(s.C(category.FieldName), *filter.Name))
		}
	}).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", repository.ErrFindCategory, err.Error())
	}
	var categories = make([]model.Category, 0, len(entcategories))
	for _, entcategory := range entcategories {
		categories = append(categories, c.ent2model(entcategory))
	}
	return categories, nil
}

func (c *Category) ent2model(category *ent.Category) model.Category {
	return model.Category{
		Id: uint(category.ID),
		Data: model.CategoryData{
			Name:        category.Name,
			Title:       category.Title,
			Description: &category.Description,
			ParentId:    (*uint)(unsafe.Pointer(category.ParentID)),
		},
	}
}
