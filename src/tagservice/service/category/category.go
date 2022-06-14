package category

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/tagservice/tagservice"
	"github.com/dmalykh/tagservice/tagservice/model"
	"github.com/dmalykh/tagservice/tagservice/repository"
	"go.uber.org/zap"
)

type Config struct {
	CategoryRepository repository.Category
	TagService         tagservice.Tag
	Logger             *zap.Logger
}

func New(config *Config) tagservice.Category {
	return &CategoryService{
		tagService:         config.TagService,
		categoryRepository: config.CategoryRepository,
		log:                config.Logger,
	}
}

type CategoryService struct {
	tagService         tagservice.Tag
	categoryRepository repository.Category
	log                *zap.Logger
}

func (c *CategoryService) Create(ctx context.Context, data *model.CategoryData) (model.Category, error) {
	var logger = c.log.With(zap.String(`method`, `Create`), zap.Any(`data`, *data))
	// Check parent's category exists
	if data.ParentId != nil {
		if _, err := c.categoryRepository.GetById(ctx, *data.ParentId); err != nil {
			logger.Error(`get parent category by id`, zap.Error(err), zap.Uintp(`parentId`, data.ParentId))
			if errors.Is(err, repository.ErrFindCategory) {
				return model.Category{}, fmt.Errorf(`parent id error: %w %d`, tagservice.ErrCategoryNotFound, *data.ParentId)
			}
			return model.Category{}, fmt.Errorf(`unknown parent id error %w`, err)
		}
	}
	// Create category
	category, err := c.categoryRepository.Create(ctx, data)
	logger.Debug(`category created`, zap.Error(err))
	if err != nil {
		return model.Category{}, fmt.Errorf(`%w %s`, tagservice.ErrCategoryNotCreated, err.Error())
	}
	return category, nil
}

func (c *CategoryService) Update(ctx context.Context, id uint, data *model.CategoryData) (model.Category, error) {
	var logger = c.log.With(zap.String(`method`, `Update`), zap.Uint("id", id))
	category, err := c.categoryRepository.GetById(ctx, id)
	if err != nil {
		logger.Error(`get category by id`, zap.Error(err))
		if errors.Is(err, repository.ErrFindCategory) {
			return model.Category{}, fmt.Errorf(`%w %d`, tagservice.ErrCategoryNotFound, id)
		}
		return model.Category{}, fmt.Errorf(`unknown error %w`, err)
	}
	// Check parent's category exists
	if data.ParentId != nil && *data.ParentId != 0 {
		if _, err := c.categoryRepository.GetById(ctx, *data.ParentId); err != nil {
			logger.Error(`get parent category by id`, zap.Error(err), zap.Uintp(`parentId`, data.ParentId))
			if errors.Is(err, repository.ErrFindCategory) {
				return model.Category{}, fmt.Errorf(`parent id error: %w %d`, tagservice.ErrCategoryNotFound, *data.ParentId)
			}
			return model.Category{}, fmt.Errorf(`unknown parent id error %w`, err)
		}
	}
	// Avoid empty values
	if data.Name == `` {
		data.Name = category.Data.Name
	}
	if data.Title == `` {
		data.Title = category.Data.Title
	}
	if data.Description == nil {
		data.Description = category.Data.Description
	}
	if data.ParentId == nil {
		data.ParentId = category.Data.ParentId
	} else if *data.ParentId == 0 {
		data.ParentId = nil
	}
	// @TODO Avoid loops with ParentId!
	category, err = c.categoryRepository.Update(ctx, category.Id, data)
	logger.Debug(`category updated`, zap.Error(err))
	if err != nil {
		return model.Category{}, fmt.Errorf(`%w %s`, tagservice.ErrCategoryNotUpdated, err.Error())
	}
	return category, nil
}

// Delete category and it's dependencies
func (c *CategoryService) Delete(ctx context.Context, id uint) error {
	var logger = c.log.With(zap.String(`method`, `Delete`), zap.Uint("id", id))
	// Check category exists
	if _, err := c.GetById(ctx, id); err != nil {
		return err
	}

	// Check tags. Category should be empty before deletion
	tags, err := c.tagService.GetList(ctx, id, 1, 0)
	logger.Debug(`get tags of category`, zap.Error(err))
	if err != nil {
		return fmt.Errorf(`unknown error %w`, err)
	}
	if len(tags) > 0 {
		return tagservice.ErrCategoryHasTags
	}

	// Delete category
	logger.Debug(`delete category`, zap.Uint(`id`, id))
	if err := c.categoryRepository.DeleteById(ctx, id); err != nil {
		return fmt.Errorf(`can't remove category %w`, err)
	}
	return nil
}

func (c *CategoryService) GetList(ctx context.Context, limit, offset uint) ([]model.Category, error) {
	var logger = c.log.With(zap.String(`method`, `GetList`), zap.Uint(`limit`, limit), zap.Uint(`offset`, offset))
	list, err := c.categoryRepository.GetList(ctx, limit, offset)
	logger.Debug(`get list`, zap.Error(err))
	if err != nil {
		return nil, fmt.Errorf(`can't receive list of categorys %w`, err)
	}
	return list, nil
}

func (c *CategoryService) GetById(ctx context.Context, id uint) (model.Category, error) {
	var logger = c.log.With(zap.String(`method`, `GetById`), zap.Uint("id", id))
	category, err := c.categoryRepository.GetById(ctx, id)
	if err != nil {
		logger.Error(`get category by id`, zap.Error(err))
		if errors.Is(err, repository.ErrFindCategory) {
			return category, fmt.Errorf(`%w %d`, tagservice.ErrCategoryNotFound, id)
		}
		return category, fmt.Errorf(`unknown error %w`, err)
	}
	return category, nil
}
