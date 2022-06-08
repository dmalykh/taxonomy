package tag

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"tagservice/server"
	"tagservice/server/model"
	"tagservice/server/repository"
	"tagservice/server/repository/transaction"
)

type Config struct {
	Transaction        transaction.Transactioner
	TagRepository      repository.Tag
	RelationRepository repository.Relation
	CategoryRepository repository.Category
	NamespaceService   server.Namespace
	Logger             *zap.Logger
}

func New(config *Config) server.Tag {
	return &TagService{
		transaction:        config.Transaction,
		relationRepository: config.RelationRepository,
		categoryRepository: config.CategoryRepository,
		namespaceService:   config.NamespaceService,
		tagRepository:      config.TagRepository,
		log:                config.Logger,
	}
}

type TagService struct {
	transaction        transaction.Transactioner
	relationRepository repository.Relation
	categoryRepository repository.Category
	namespaceService   server.Namespace
	tagRepository      repository.Tag
	log                *zap.Logger
}

var ErrTagNotFound = errors.New(`tag not found`)
var ErrTagNamespaceNotFound = errors.New(`tag's namespace not found`)
var ErrTagNotCreated = errors.New(`tag had not created`)
var ErrTagRelationNotCreated = errors.New(`tag's relation had not created`)
var ErrTagNotUpdated = errors.New(`tag have not updated`)

func (t *TagService) Create(ctx context.Context, data *model.TagData) (model.Tag, error) {
	var logger = t.log.With(zap.String(`method`, `Create`), zap.Any(`data`, *data))

	// Check category exists
	if _, err := t.categoryRepository.GetById(ctx, data.CategoryId); err != nil {
		logger.Error(`get category by id`, zap.Error(err), zap.Uint(`categoryId`, data.CategoryId))
		if errors.Is(err, repository.ErrFindCategory) {
			return model.Tag{}, fmt.Errorf(`%w %d`, server.ErrCategoryNotFound, data.CategoryId)
		}
		return model.Tag{}, fmt.Errorf(`unknown category error %w`, err)
	}

	tag, err := t.tagRepository.Create(ctx, data)
	logger.Debug(`tag created`, zap.Any(`tag`, tag), zap.Error(err))
	if err != nil {
		return model.Tag{}, fmt.Errorf(`%w %s`, ErrTagNotCreated, err.Error())
	}
	return tag, nil
}

func (t *TagService) Update(ctx context.Context, id uint, data *model.TagData) (model.Tag, error) {
	var logger = t.log.With(zap.String(`method`, `Update`), zap.Uint("id", id), zap.Any(`data`, *data))

	// Check tag exists
	tag, err := t.tagRepository.GetById(ctx, id)
	if err != nil {
		logger.Error(`get tag by id`, zap.Error(err))
		if errors.Is(err, repository.ErrFindTag) {
			return model.Tag{}, fmt.Errorf(`%w %d`, ErrTagNotFound, id)
		}
		return model.Tag{}, fmt.Errorf(`unknown error %w`, err)
	}
	// Avoid empty values
	if data.Name == `` {
		data.Name = tag.Data.Name
	}
	if data.Title == `` {
		data.Title = tag.Data.Title
	}
	if data.Description == `` {
		data.Description = tag.Data.Description
	}
	if data.CategoryId == 0 {
		data.CategoryId = tag.Data.CategoryId
	}
	// Check category exists
	if _, err := t.categoryRepository.GetById(ctx, data.CategoryId); err != nil {
		logger.Error(`get category by id`, zap.Error(err), zap.Uint(`categoryId`, data.CategoryId))
		if errors.Is(err, repository.ErrFindCategory) {
			return model.Tag{}, fmt.Errorf(`%w %d`, server.ErrCategoryNotFound, data.CategoryId)
		}
		return model.Tag{}, fmt.Errorf(`unknown category error %w`, err)
	}
	// Update tag
	updated, err := t.tagRepository.Update(ctx, tag.Id, data)
	logger.Debug(`tag updated`, zap.Any(`tag`, updated), zap.Error(err))
	if err != nil {
		return model.Tag{}, fmt.Errorf(`%w %s`, ErrTagNotUpdated, err.Error())
	}
	return updated, nil
}

func (t *TagService) Delete(ctx context.Context, id uint) error {
	var logger = t.log.With(zap.String(`method`, `Delete`), zap.Uint("id", id))

	// Check tag exists
	tag, err := t.tagRepository.GetById(ctx, id)
	if err != nil {
		logger.Error(`get tag by id`, zap.Error(err))
		if errors.Is(err, repository.ErrFindTag) {
			return fmt.Errorf(`%w %d`, ErrTagNotFound, id)
		}
		return fmt.Errorf(`unknown error %w`, err)
	}

	tx, err := t.transaction.BeginTx(ctx)
	logger.Debug(`start Transaction`, zap.Error(err))
	if err != nil {
		return fmt.Errorf(`transaction error %w`, err)
	}

	// Delete relations with this tag
	logger.Debug(`delete relations by tag id`, zap.Uint(`id`, tag.Id))
	if err := tx.Relation().Delete(ctx, []uint{tag.Id}, nil, nil); err != nil {
		logger.Error(`rollback`, zap.Error(err))
		if err := tx.Rollback(ctx); err != nil {
			return fmt.Errorf(`rollback error %w`, err)
		}
		return fmt.Errorf(`can't remove relations %w`, err)
	}
	// Delete tag
	logger.Debug(`delete tag by id`, zap.Uint(`id`, tag.Id))
	if err := tx.Tag().DeleteById(ctx, tag.Id); err != nil {
		logger.Error(`rollback`, zap.Error(err))
		if err := tx.Rollback(ctx); err != nil {
			return fmt.Errorf(`rollback error %w`, err)
		}
		return fmt.Errorf(`can't remove tag %w`, err)
	}

	logger.Debug(`commit`)
	if err := tx.Commit(ctx); err != nil {
		logger.Error(`not committed`, zap.Error(err))
		return fmt.Errorf(`commit error %w`, err)
	}
	return nil
}

func (t *TagService) GetById(ctx context.Context, id uint) (model.Tag, error) {
	var logger = t.log.With(zap.String(`method`, `GetById`), zap.Uint("id", id))

	tag, err := t.tagRepository.GetById(ctx, id)
	if err != nil {
		logger.Error(`get tag by id`, zap.Error(err))
		if errors.Is(err, repository.ErrFindTag) {
			return tag, fmt.Errorf(`%w %d`, ErrTagNotFound, id)
		}
		return tag, fmt.Errorf(`unknown error %w`, err)
	}
	return tag, nil
}

func (t *TagService) GetByName(ctx context.Context, name string, categoryId uint) (model.Tag, error) {
	var logger = t.log.With(zap.String(`method`, `GetByName`), zap.String("name", name))

	// Check category exists
	if _, err := t.categoryRepository.GetById(ctx, categoryId); err != nil {
		logger.Error(`get category by id`, zap.Error(err), zap.Uint(`categoryId`, categoryId))
		if errors.Is(err, repository.ErrFindCategory) {
			return model.Tag{}, fmt.Errorf(`%w %d`, server.ErrCategoryNotFound, categoryId)
		}
		return model.Tag{}, fmt.Errorf(`unknown category error %w`, err)
	}

	tags, err := t.tagRepository.GetByName(ctx, name)
	if err != nil {
		logger.Error(`get tag by name`, zap.Error(err))
		if errors.Is(err, repository.ErrFindTag) {
			return model.Tag{}, fmt.Errorf(`%w %s`, ErrTagNotFound, name)
		}
		return model.Tag{}, fmt.Errorf(`unknown error %w`, err)
	}
	for _, tag := range tags {
		if tag.Data.CategoryId == categoryId {
			return tag, nil
		}
	}
	return model.Tag{}, fmt.Errorf(`%w with %q, %d`, ErrTagNotFound, name, categoryId)
}

func (t *TagService) SetRelation(ctx context.Context, tagId uint, entitiesNamespace string, entitiesId ...uint) error {
	namespace, err := t.namespaceService.GetByName(ctx, entitiesNamespace)
	if err != nil {
		return fmt.Errorf(`%w %s`, ErrTagNamespaceNotFound, err.Error())
	}
	var relations = make([]*model.Relation, 0, len(entitiesId))
	for _, entityId := range entitiesId {
		relations = append(relations, &model.Relation{
			TagId:       tagId,
			NamespaceId: namespace.Id,
			EntityId:    entityId,
		})
	}
	if err := t.relationRepository.Create(ctx, relations...); err != nil {
		return fmt.Errorf(`%w %s`, ErrTagRelationNotCreated, err.Error())
	}
	return nil
}

func (t *TagService) GetList(ctx context.Context, categoryId uint, limit, offset uint) ([]model.Tag, error) {
	var logger = t.log.With(zap.String(`method`, `GetList`), zap.Uint(`categoryId`, categoryId), zap.Uint(`limit`, limit), zap.Uint(`offset`, offset))

	tags, err := t.tagRepository.GetByFilter(ctx, model.TagFilter{CategoryId: []uint{categoryId}}, limit, offset)
	logger.Debug(`got tags`, zap.Any(`tags`, tags), zap.Error(err))
	if err != nil {
		return nil, fmt.Errorf(`unknown error %w`, err)
	}
	return tags, nil
}

func (t *TagService) GetRelationEntities(ctx context.Context, namespaceName string, tagGroups [][]uint) ([]model.Relation, error) {
	var logger = t.log.With(zap.String(`method`, `GetRelationEntities`), zap.String(`namespaceName`, namespaceName), zap.Any(`tagGroups`, tagGroups))

	namespace, err := t.namespaceService.GetByName(ctx, namespaceName)
	logger.Debug(`got namespace`, zap.Any(`namespace`, namespace), zap.Error(err))
	if err != nil {
		return nil, ErrTagNamespaceNotFound
	}

	var unique = make(map[uint]model.Relation)
	for _, tagIds := range tagGroups {
		rels, err := t.relationRepository.Get(ctx, tagIds, []uint{namespace.Id}, nil)
		logger.Debug(`got relations`, zap.Uints(`tagIds`, tagIds), zap.Any(`rels`, rels), zap.Error(err))
		if err != nil {
			return nil, fmt.Errorf(`unknown error %w`, err)
		}
		for _, rel := range rels {
			unique[rel.EntityId] = rel
		}
	}
	logger.Debug(`finally`, zap.Any(`unique`, unique))

	var relations = make([]model.Relation, 0, len(unique))
	for _, rel := range unique {
		relations = append(relations, rel)
	}
	return relations, nil
}

func (t *TagService) GetTagsByEntities(ctx context.Context, namespaceName string, entities ...uint) ([]model.Tag, error) {
	var logger = t.log.With(zap.String(`method`, `GetTagsByEntities`), zap.String(`namespaceName`, namespaceName), zap.Uints(`entities`, entities))

	namespace, err := t.namespaceService.GetByName(ctx, namespaceName)
	logger.Debug(`got namespace`, zap.Any(`namespace`, namespace), zap.Error(err))
	if err != nil {
		return nil, ErrTagNamespaceNotFound
	}

	relations, err := t.relationRepository.Get(ctx, nil, []uint{namespace.Id}, entities)
	logger.Debug(`got relations`, zap.Any(`relations`, relations), zap.Error(err))
	if err != nil {
		return nil, fmt.Errorf(`unknown error %w`, err)
	}

	var tags = make([]model.Tag, 0, len(relations))
	for _, relation := range relations {
		// Change for one request if there is a lot of relation will be found
		tag, err := t.tagRepository.GetById(ctx, relation.TagId)
		logger.Debug(`got tag`, zap.Uint(`id`, relation.TagId), zap.Any(`tag`, tag), zap.Error(err))
		if err != nil {
			logger.DPanic(`unknown tag in relation`, zap.Error(err))
			continue
		}
		tags = append(tags, tag)
	}
	return tags, nil
}
