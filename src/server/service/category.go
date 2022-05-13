package service

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"tagservice/server"
	"tagservice/server/model"
	"tagservice/server/repository"
)

type CategoryService struct {
	tagService         server.Tag
	categoryRepository repository.Category
	log                *zap.Logger
}

var ErrCategoryNotFound = errors.New(`category not found`)
var ErrCategoryNotCreated = errors.New(`category has not created`)
var ErrCategoryHasTags = errors.New(`category has tags, but should be empty`)
var ErrCategoryNotUpdated = errors.New(`category has not updated`)

func (c *CategoryService) Create(ctx context.Context, data *model.CategoryData) (model.Category, error) {
	var logger = c.log.With(zap.String(`method`, `Create`), zap.Any(`data`, *data))
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	category, err := c.categoryRepository.Create(ctx, data)
	logger.Debug(`category created`, zap.Error(err))
	if err != nil {
		return model.Category{}, fmt.Errorf(`%w %s`, ErrCategoryNotCreated, err.Error())
	}
	return category, nil
}

func (c *CategoryService) Update(ctx context.Context, id uint64, data *model.CategoryData) (model.Category, error) {
	var logger = c.log.With(zap.String(`method`, `Update`), zap.Uint64("id", id))
	category, err := c.categoryRepository.GetById(ctx, id)
	if err != nil {
		logger.Error(`get category by id`, zap.Error(err))
		if errors.Is(err, repository.ErrNotFound) {
			return model.Category{}, fmt.Errorf(`%w %d`, ErrCategoryNotFound, id)
		}
		return model.Category{}, fmt.Errorf(`unknown error %w`, err)
	}
	category, err = c.categoryRepository.Update(ctx, category.Id, data)
	logger.Debug(`category updated`, zap.Error(err))
	if err != nil {
		return model.Category{}, fmt.Errorf(`%w %s`, ErrCategoryNotUpdated, err.Error())
	}
	return category, nil
}

// Delete category and it's dependencies
func (c *CategoryService) Delete(ctx context.Context, id uint64) error {
	var logger = c.log.With(zap.String(`method`, `Delete`), zap.Uint64("id", id))
	// Check category exists
	if _, err := c.GetById(ctx, id); err != nil {
		return err
	}

	// Check tags. Category should be empty before deletion
	tags, err := c.tagService.GetList(ctx, nil, id, 1, 0)
	logger.Debug(`get tags of category`, zap.Error(err))
	if err != nil {
		return fmt.Errorf(`unknown error %w`, err)
	}
	if len(tags) > 0 {
		return ErrCategoryHasTags
	}

	// Delete category
	logger.Debug(`delete category`, zap.Uint64(`id`, id))
	if err := c.categoryRepository.DeleteById(ctx, id); err != nil {
		return fmt.Errorf(`can't remove category %w`, err)
	}
	return nil
}

func (c *CategoryService) GetList(ctx context.Context, limit, offset uint64) ([]model.Category, error) {
	var logger = c.log.With(zap.String(`method`, `GetList`), zap.Uint64(`limit`, limit), zap.Uint64(`offset`, offset))
	list, err := c.categoryRepository.GetList(ctx, limit, offset)
	logger.Debug(`get list`, zap.Error(err))
	if err != nil {
		return nil, fmt.Errorf(`can't receive list of categorys %w`, err)
	}
	return list, nil
}

func (c *CategoryService) GetById(ctx context.Context, id uint64) (model.Category, error) {
	var logger = c.log.With(zap.String(`method`, `GetById`), zap.Uint64("id", id))
	category, err := c.categoryRepository.GetById(ctx, id)
	if err != nil {
		logger.Error(`get category by id`, zap.Error(err))
		if errors.Is(err, repository.ErrNotFound) {
			return category, fmt.Errorf(`%w %d`, ErrCategoryNotFound, id)
		}
		return category, fmt.Errorf(`unknown error %w`, err)
	}
	return category, nil
}
