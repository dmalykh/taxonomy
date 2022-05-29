package repository

import (
	"context"
	"errors"
	"fmt"
	"tagservice/repository/entgo/ent"
	"tagservice/server/model"
	"tagservice/server/repository"
	"unsafe"
)

var (
	ErrCreateCategory = errors.New(`failed to create category`)
	ErrUpdateCategory = errors.New(`failed to update category`)
	ErrFindCategory   = errors.New(`failed to find category`)
	ErrDeleteCategory = errors.New(`failed to delete category`)
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
		SetDescription(data.Description).
		SetNillableParentID(func() *int { return (*int)(unsafe.Pointer(data.PatentId)) }()).
		Save(ctx)

	if err != nil {
		return model.Category{}, fmt.Errorf("%w: %s", ErrCreateCategory, err.Error())
	}
	return c.ent2model(ns), nil
}

func (c *Category) Update(ctx context.Context, id uint, data *model.CategoryData) (model.Category, error) {
	category, err := c.client.UpdateOneID(int(id)).
		SetName(data.Name).
		SetTitle(data.Title).
		SetDescription(data.Description).
		SetNillableParentID(func() *int { return (*int)(unsafe.Pointer(data.PatentId)) }()).
		Save(ctx)
	if err != nil {
		return model.Category{}, fmt.Errorf("%w: %s", ErrUpdateCategory, err.Error())
	}
	return c.ent2model(category), err
}

func (c *Category) DeleteById(ctx context.Context, id uint) error {
	if err := c.client.DeleteOneID(int(id)).Exec(ctx); err != nil {
		return fmt.Errorf("%w (%d): %s", ErrDeleteCategory, id, err.Error())
	}
	return nil
}

func (c *Category) GetById(ctx context.Context, id uint) (model.Category, error) {
	ns, err := c.client.Get(ctx, int(id))
	if err != nil {
		return model.Category{}, fmt.Errorf("%w (%d): %s", ErrFindCategory, id, err.Error())
	}
	return c.ent2model(ns), err
}

func (c *Category) GetList(ctx context.Context, limit, offset uint) ([]model.Category, error) {
	entcategories, err := c.client.Query().Offset(int(offset)).Limit(int(limit)).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFindCategory, err.Error())
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
			Description: category.Description,
			PatentId:    (*uint)(unsafe.Pointer(category.ParentID)),
		},
	}
}