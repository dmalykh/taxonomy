package namespace

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
	Transaction         transaction.Transactioner
	NamespaceRepository repository.Namespace
	RelationRepository  repository.Relation
	Logger              *zap.Logger
}

func New(config *Config) tagservice.Namespace {
	return &NamespaceService{
		transaction:         config.Transaction,
		namespaceRepository: config.NamespaceRepository,
		relationRepository:  config.RelationRepository,
		log:                 config.Logger,
	}
}

//goland:noinspection GoNameStartsWithPackageName,GoUnnecessarilyExportedIdentifiers
//nolint:revive
type NamespaceService struct {
	transaction         transaction.Transactioner
	namespaceRepository repository.Namespace
	relationRepository  repository.Relation
	log                 *zap.Logger
}

func (n *NamespaceService) Create(ctx context.Context, name string) (model.Namespace, error) {
	logger := n.log.With(zap.String(`method`, `Create`), zap.String(`name`, name))
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
		}
	}(logger)

	namespace, err := n.namespaceRepository.Create(ctx, name)
	logger.Debug(`namespace created`, zap.Error(err))

	if err != nil {
		return model.Namespace{}, fmt.Errorf(`%w %s`, tagservice.ErrNamespaceNotCreated, err.Error())
	}

	return namespace, nil
}

func (n *NamespaceService) Update(ctx context.Context, id uint, name string) (model.Namespace, error) {
	logger := n.log.With(zap.String(`method`, `Update`), zap.Uint("id", id))
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	namespace, err := n.namespaceRepository.GetByID(ctx, id)
	if err != nil {
		logger.Error(`get namespace by id`, zap.Error(err))

		if errors.Is(err, repository.ErrFindNamespace) {
			return model.Namespace{}, fmt.Errorf(`%w %d`, tagservice.ErrNamespaceNotFound, id)
		}

		return model.Namespace{}, fmt.Errorf(`unknown error %w`, err)
	}

	namespace, err = n.namespaceRepository.Update(ctx, namespace.ID, name)
	logger.Debug(`namespace updated`, zap.Error(err))

	if err != nil {
		return model.Namespace{}, fmt.Errorf(`%w %s`, tagservice.ErrNamespaceNotUpdated, err.Error())
	}

	return namespace, nil
}

// Delete namespace and it's dependencies.
func (n *NamespaceService) Delete(ctx context.Context, id uint) error {
	logger := n.log.With(zap.String(`method`, `Delete`), zap.Uint("id", id))
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	namespace, err := n.namespaceRepository.GetByID(ctx, id)
	if err != nil {
		logger.Error(`get namespace by id`, zap.Error(err))

		if errors.Is(err, repository.ErrFindNamespace) {
			return fmt.Errorf(`%w %d`, tagservice.ErrNamespaceNotFound, id)
		}

		return fmt.Errorf(`unknown error %w`, err)
	}

	tx, err := n.transaction.BeginTx(ctx)
	logger.Debug(`start transaction`, zap.Error(err))

	if err != nil {
		return fmt.Errorf(`transaction error %w`, err)
	}

	// Delete dependencies
	logger.Debug(`delete namespace by id`, zap.Uint(`id`, namespace.ID))

	if err := tx.Relation().Delete(ctx, nil, []uint{namespace.ID}, nil); err != nil {
		logger.Error(`rollback`, zap.Error(err))

		if err := tx.Rollback(ctx); err != nil {
			return fmt.Errorf(`rollback error %w`, err)
		}

		return fmt.Errorf(`can't remove relations %w`, err)
	}

	// Delete namespace
	logger.Debug(`delete namespace by id`, zap.Uint(`id`, namespace.ID))

	if err := tx.Namespace().DeleteByID(ctx, namespace.ID); err != nil {
		logger.Error(`rollback`, zap.Error(err))

		if err := tx.Rollback(ctx); err != nil {
			return fmt.Errorf(`rollback error %w`, err)
		}

		return fmt.Errorf(`can't remove namespace %w`, err)
	}

	logger.Debug(`commit`)

	if err := tx.Commit(ctx); err != nil {
		logger.Error(`not committed`, zap.Error(err))

		return fmt.Errorf(`commit error %w`, err)
	}

	return nil
}

func (n *NamespaceService) GetList(ctx context.Context, limit, offset uint) ([]model.Namespace, error) {
	logger := n.log.With(zap.String(`method`, `GetList`), zap.Uint(`limit`, limit), zap.Uint(`offset`, offset))
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	list, err := n.namespaceRepository.GetList(ctx, limit, offset)
	logger.Debug(`get list`, zap.Error(err))

	if err != nil {
		return nil, fmt.Errorf(`can't receive list of namespaces %w`, err)
	}

	return list, nil
}

func (n *NamespaceService) GetByName(ctx context.Context, name string) (model.Namespace, error) {
	logger := n.log.With(zap.String(`method`, `GetByName`), zap.String("name", name))
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	namespace, err := n.namespaceRepository.GetByName(ctx, name)
	if err != nil {
		logger.Error(`get namespace by name`, zap.Error(err))

		if errors.Is(err, repository.ErrFindNamespace) {
			return namespace, fmt.Errorf(`%w %s`, tagservice.ErrNamespaceNotFound, name)
		}

		return namespace, fmt.Errorf(`unknown error %w`, err)
	}

	return namespace, nil
}
