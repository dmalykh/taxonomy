package tag

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/tagservice/tagservice"
	"github.com/dmalykh/tagservice/tagservice/model"
	"github.com/dmalykh/tagservice/tagservice/repository"
	"github.com/dmalykh/tagservice/tagservice/repository/transaction"
	"go.uber.org/zap"
)

type Config struct {
	Transaction        transaction.Transactioner
	TagRepository      repository.Tag
	RelationRepository repository.Relation
	CategoryRepository repository.Category
	NamespaceService   tagservice.Namespace
	Logger             *zap.Logger
}

func New(config *Config) tagservice.Tag {
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
	namespaceService   tagservice.Namespace
	tagRepository      repository.Tag
	log                *zap.Logger
}

func (t *TagService) Create(ctx context.Context, data *model.TagData) (model.Tag, error) {
	var logger = t.log.With(zap.String(`method`, `Create`), zap.Any(`data`, *data))

	// Check category exists
	if _, err := t.categoryRepository.GetById(ctx, data.CategoryId); err != nil {
		logger.Error(`get category by id`, zap.Error(err), zap.Uint(`categoryId`, data.CategoryId))
		if errors.Is(err, repository.ErrFindCategory) {
			return model.Tag{}, fmt.Errorf(`%w %d`, tagservice.ErrCategoryNotFound, data.CategoryId)
		}
		return model.Tag{}, fmt.Errorf(`unknown category error %w`, err)
	}

	tag, err := t.tagRepository.Create(ctx, data)
	logger.Debug(`tag created`, zap.Any(`tag`, tag), zap.Error(err))
	if err != nil {
		return model.Tag{}, fmt.Errorf(`%w %s`, tagservice.ErrTagNotCreated, err.Error())
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
			return model.Tag{}, fmt.Errorf(`%w %d`, tagservice.ErrTagNotFound, id)
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
			return model.Tag{}, fmt.Errorf(`%w %d`, tagservice.ErrCategoryNotFound, data.CategoryId)
		}
		return model.Tag{}, fmt.Errorf(`unknown category error %w`, err)
	}
	// Update tag
	updated, err := t.tagRepository.Update(ctx, tag.Id, data)
	logger.Debug(`tag updated`, zap.Any(`tag`, updated), zap.Error(err))
	if err != nil {
		return model.Tag{}, fmt.Errorf(`%w %s`, tagservice.ErrTagNotUpdated, err.Error())
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
			return fmt.Errorf(`%w %d`, tagservice.ErrTagNotFound, id)
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
			return tag, fmt.Errorf(`%w %d`, tagservice.ErrTagNotFound, id)
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
			return model.Tag{}, fmt.Errorf(`%w %d`, tagservice.ErrCategoryNotFound, categoryId)
		}
		return model.Tag{}, fmt.Errorf(`unknown category error %w`, err)
	}

	tags, err := t.tagRepository.GetByName(ctx, name)
	if err != nil {
		logger.Error(`get tag by name`, zap.Error(err))
		if errors.Is(err, repository.ErrFindTag) {
			return model.Tag{}, fmt.Errorf(`%w %s`, tagservice.ErrTagNotFound, name)
		}
		return model.Tag{}, fmt.Errorf(`unknown error %w`, err)
	}
	for _, tag := range tags {
		if tag.Data.CategoryId == categoryId {
			return tag, nil
		}
	}
	return model.Tag{}, fmt.Errorf(`%w with %q, %d`, tagservice.ErrTagNotFound, name, categoryId)
}

func (t *TagService) SetRelation(ctx context.Context, tagId uint, entitiesNamespace string, entitiesId ...uint) error {
	var logger = t.log.With(zap.String(`method`, `SetRelation`), zap.Uint(`tagId`, tagId), zap.String("entitiesNamespace", entitiesNamespace), zap.Uints(`entitiesId`, entitiesId))

	namespace, err := t.namespaceService.GetByName(ctx, entitiesNamespace)
	if err != nil {
		return fmt.Errorf(`%w %s`, tagservice.ErrTagNamespaceNotFound, err.Error())
	}

	// Check tag exists
	if _, err := t.tagRepository.GetById(ctx, tagId); err != nil {
		logger.Error(`get tag by id`, zap.Error(err), zap.Uint(`tagId`, tagId))
		if errors.Is(err, repository.ErrFindTag) {
			return fmt.Errorf(`%w %d`, tagservice.ErrTagNotFound, tagId)
		}
		return fmt.Errorf(`unknown tag error %w`, err)
	}

	// Prepare relations without duplicates
	var relations = make([]*model.Relation, 0, len(entitiesId))
	var seen = make(map[uint]struct{})
	for _, entityId := range entitiesId {
		if _, exists := seen[entityId]; exists {
			continue
		}
		relations = append(relations, &model.Relation{
			TagId:       tagId,
			NamespaceId: namespace.Id,
			EntityId:    entityId,
		})
		seen[entityId] = struct{}{}
	}

	// Delete exists and insert prepared
	tx, err := t.transaction.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf(`transaction error %w`, err)
	}

	if err := tx.Relation().Delete(ctx, []uint{tagId}, []uint{namespace.Id}, entitiesId); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return fmt.Errorf(`rollback error %w`, err)
		}
		return fmt.Errorf(`can't remove relation %w: %s`, tagservice.ErrTagRelationNotRemoved, err.Error())
	}

	if err := tx.Relation().Create(ctx, relations...); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return fmt.Errorf(`rollback error %w`, err)
		}
		return fmt.Errorf(`can't create relation %w: %s`, tagservice.ErrTagRelationNotCreated, err.Error())
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf(`commit error %w`, err)
	}
	return nil
}

func (t *TagService) UnsetRelation(ctx context.Context, tagId uint, entitiesNamespace string, entitiesId ...uint) error {
	var logger = t.log.With(zap.String(`method`, `UnsetRelation`), zap.Uint(`tagId`, tagId), zap.String("entitiesNamespace", entitiesNamespace), zap.Uints(`entitiesId`, entitiesId))

	namespace, err := t.namespaceService.GetByName(ctx, entitiesNamespace)
	if err != nil {
		return fmt.Errorf(`%w %s`, tagservice.ErrTagNamespaceNotFound, err.Error())
	}

	// Check tag exists
	if _, err := t.tagRepository.GetById(ctx, tagId); err != nil {
		logger.Error(`get tag by id`, zap.Error(err), zap.Uint(`tagId`, tagId))
		if errors.Is(err, repository.ErrFindTag) {
			return fmt.Errorf(`%w %d`, tagservice.ErrTagNotFound, tagId)
		}
		return fmt.Errorf(`unknown tag error %w`, err)
	}
	// Remove relations
	if err := t.relationRepository.Delete(ctx, []uint{tagId}, []uint{namespace.Id}, entitiesId); err != nil {
		return fmt.Errorf(`can't remove relation %w: %s`, tagservice.ErrTagRelationNotRemoved, err.Error())
	}

	return nil
}

func (t *TagService) GetList(ctx context.Context, filter *model.TagFilter) ([]model.Tag, error) {
	var logger = t.log.With(zap.String(`method`, `GetList`), zap.Any(`filter`, filter))

	tags, err := t.tagRepository.GetByFilter(ctx, filter)
	logger.Debug(`got tags`, zap.Any(`tags`, tags), zap.Error(err))
	if err != nil {
		return nil, fmt.Errorf(`unknown error %w`, err)
	}
	return tags, nil
}

func (t *TagService) GetTagsByEntities(ctx context.Context, namespaceName string, entities ...uint) ([]model.Tag, error) {
	var logger = t.log.With(zap.String(`method`, `GetTagsByEntities`), zap.String(`namespaceName`, namespaceName), zap.Uints(`entities`, entities))

	namespace, err := t.namespaceService.GetByName(ctx, namespaceName)
	logger.Debug(`got namespace`, zap.Any(`namespace`, namespace), zap.Error(err))
	if err != nil {
		return nil, tagservice.ErrTagNamespaceNotFound
	}

	relations, err := t.relationRepository.Get(ctx, &model.RelationFilter{
		Namespace: []uint{namespace.Id},
		EntityId:  entities,
	})
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

func (t *TagService) GetRelations(ctx context.Context, filter *model.EntityFilter) ([]model.Relation, error) {
	var logger = t.log.With(zap.String(`method`, `GetTagsByEntities`), zap.Any(`filter`, filter))

	// Get ids of namespaces
	var namespaces = make([]uint, 0, len(filter.Namespace))
	for _, ns := range filter.Namespace {
		namespace, err := t.namespaceService.GetByName(ctx, ns)
		if err != nil {
			return nil, tagservice.ErrTagNamespaceNotFound
		}
		namespaces = append(namespaces, namespace.Id)
	}

	// Get relations
	relations, err := t.relationRepository.Get(ctx, &model.RelationFilter{
		TagId:     filter.TagId,
		EntityId:  filter.EntityId,
		Namespace: namespaces,
		AfterId:   filter.AfterId,
		Limit:     filter.Limit,
	})
	logger.Debug(`got relations`, zap.Any(`relations`, relations), zap.Error(err))
	if err != nil {
		return nil, fmt.Errorf(`unknown error %w`, err)
	}
	return relations, nil
}
