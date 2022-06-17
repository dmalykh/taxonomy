package repository

import (
	"context"
	"fmt"
	"unsafe"

	"entgo.io/ent/dialect/sql"
	"github.com/dmalykh/tagservice/repository/entgo/ent"
	"github.com/dmalykh/tagservice/repository/entgo/ent/category"
	"github.com/dmalykh/tagservice/tagservice/model"
	"github.com/dmalykh/tagservice/tagservice/repository"
)

//goland:noinspection GoUnnecessarilyExportedIdentifiers
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
		SetNillableParentID(func() *int { return (*int)(unsafe.Pointer(data.ParentID)) }()).
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
	updated, err := c.client.UpdateOneID(int(id)).
		SetName(data.Name).
		SetTitle(data.Title).
		SetNillableDescription(data.Description).
		ClearParentID().
		SetNillableParentID(func() *int { return (*int)(unsafe.Pointer(data.ParentID)) }()).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return model.Category{}, repository.ErrNotUniqueName
		}

		return model.Category{}, fmt.Errorf("%w: %s", repository.ErrUpdateCategory, err.Error())
	}

	return c.ent2model(updated), err
}

func (c *Category) DeleteByID(ctx context.Context, id uint) error {
	if err := c.client.DeleteOneID(int(id)).Exec(ctx); err != nil {
		return fmt.Errorf("%w (%d): %s", repository.ErrDeleteCategory, id, err.Error())
	}

	return nil
}

func (c *Category) GetByID(ctx context.Context, id uint) (model.Category, error) {
	ns, err := c.client.Get(ctx, int(id))
	if err != nil {
		return model.Category{}, fmt.Errorf("%w (%d): %s", repository.ErrFindCategory, id, err.Error())
	}

	return c.ent2model(ns), err
}

func (c *Category) GetList(ctx context.Context, filter *model.CategoryFilter) ([]model.Category, error) {
	entcategories, err := c.client.Query().Where(func(s *sql.Selector) {
		// Filter by parent id
		if filter.ParentID != nil {
			s.Where(sql.EQ(category.FieldParentID, *filter.ParentID))
		}
		// Filter by name
		if filter.Name != nil {
			s.Where(sql.EQ(s.C(category.FieldName), *filter.Name))
		}
	}).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", repository.ErrFindCategory, err.Error())
	}

	categories := make([]model.Category, 0, len(entcategories))

	for _, entcategory := range entcategories {
		categories = append(categories, c.ent2model(entcategory))
	}

	return categories, nil
}

func (c *Category) ent2model(category *ent.Category) model.Category {
	return model.Category{
		ID: uint(category.ID),
		Data: model.CategoryData{
			Name:        category.Name,
			Title:       category.Title,
			Description: &category.Description,
			ParentID:    (*uint)(unsafe.Pointer(category.ParentID)),
		},
	}
}
