package namespace

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/taxonomy/taxonomy"
	"github.com/dmalykh/taxonomy/taxonomy/model"
	"github.com/dmalykh/taxonomy/taxonomy/repository"
	"go.uber.org/zap"
)

type Config struct {
	NamespaceRepository repository.Namespace
	ReferenceService    taxonomy.Reference
	Logger              *zap.Logger
}

func New(config *Config) taxonomy.Namespace {
	return &NamespaceService{
		namespaceRepository: config.NamespaceRepository,
		referenceService:    config.ReferenceService,
		log:                 config.Logger,
	}
}

type NamespaceService struct {
	namespaceRepository repository.Namespace
	referenceService    taxonomy.Reference
	log                 *zap.Logger
}

func (n *NamespaceService) Create(ctx context.Context, name string) (*model.Namespace, error) {
	logger := n.log.With(zap.String(`method`, `Create`), zap.String(`name`, name))

	ns, err := n.namespaceRepository.Create(ctx, &model.NamespaceData{
		Name: name,
	})
	logger.Debug(`namespace created`, zap.Error(err))

	if err != nil {
		return nil, fmt.Errorf(`%w %s`, taxonomy.ErrNamespaceNotCreated, err.Error())
	}

	return ns, nil
}

func (n *NamespaceService) Update(ctx context.Context, id uint64, name string) (*model.Namespace, error) {
	logger := n.log.With(zap.String(`method`, `Update`), zap.Uint64("id", id))

	nss, err := n.namespaceRepository.Get(ctx, &repository.NamespaceFilter{
		ID: []uint64{id},
	})
	if err != nil {
		logger.Error(`get namespace by id`, zap.Error(err))

		if errors.Is(err, repository.ErrFindNamespace) {
			return nil, fmt.Errorf(`%w %d`, taxonomy.ErrNamespaceNotFound, id)
		}

		return nil, fmt.Errorf(`unknown error %w`, err)
	}

	if len(nss) != 1 {
		return nil, fmt.Errorf(`%w, got %d results`, taxonomy.ErrNamespaceNotFound, len(nss))
	}

	ns, err := n.namespaceRepository.Update(ctx, nss[0].ID, &model.NamespaceData{
		Name: name,
	})
	logger.Debug(`namespace updated`, zap.Error(err))

	if err != nil {
		return nil, errors.Join(taxonomy.ErrNamespaceNotUpdated, err)
	}

	return ns, nil
}

// Delete namespace and it's dependencies.
func (n *NamespaceService) Delete(ctx context.Context, id uint64) error {
	logger := n.log.With(zap.String(`method`, `Delete`), zap.Uint64("id", id))

	nss, err := n.namespaceRepository.Get(ctx, &repository.NamespaceFilter{
		ID: []uint64{id},
	})
	if err != nil {
		logger.Error(`get namespace by id`, zap.Error(err))

		if errors.Is(err, repository.ErrFindNamespace) {
			return fmt.Errorf(`%w %d`, taxonomy.ErrNamespaceNotFound, id)
		}

		return fmt.Errorf(`unknown error %w`, err)
	}

	if len(nss) != 1 {
		return fmt.Errorf(`%w, got %d results`, taxonomy.ErrNamespaceNotFound, len(nss))
	}

	// Reference exists check
	logger.Debug(`check references`)

	ref, err := n.referenceService.Get(ctx, &model.ReferenceFilter{
		Namespace: []string{nss[0].Data.Name},
	})
	if err != nil {
		logger.Error(`get references by namespace`, zap.Uint64(`namespace_id`, id), zap.Error(err))

		return fmt.Errorf(`get references by term error: %w`, err)
	}

	if len(ref) > 0 {
		return errors.Join(taxonomy.ErrReferenceExists, fmt.Errorf(`%d has %d references`, id, len(ref)))
	}

	// Delete namespace
	logger.Debug(`delete namespace by id`, zap.Uint64(`id`, nss[0].ID))

	if err := n.namespaceRepository.Delete(ctx, &repository.NamespaceFilter{
		ID: []uint64{nss[0].ID},
	}); err != nil {
		logger.Error(`delete error`, zap.Error(err))

		return errors.Join(taxonomy.ErrNamespaceNotDeleted, err)
	}

	return nil
}

func (n *NamespaceService) Get(ctx context.Context, limit uint, afterId *uint64) ([]*model.Namespace, error) {
	logger := n.log.With(zap.String(`method`, `Get`),
		zap.Uint(`limit`, limit), zap.Uint64p(`afterId`, afterId))

	list, err := n.namespaceRepository.Get(ctx, &repository.NamespaceFilter{
		Limit:   limit,
		AfterID: afterId,
	})
	logger.Debug(`get list`, zap.Error(err))

	if err != nil {
		return nil, fmt.Errorf(`can't receive list of namespaces %w`, err)
	}

	return list, nil
}

func (n *NamespaceService) GetByName(ctx context.Context, name string) (*model.Namespace, error) {
	logger := n.log.With(zap.String(`method`, `GetByName`), zap.String("name", name))

	nss, err := n.namespaceRepository.Get(ctx, &repository.NamespaceFilter{
		Name: []string{name},
	})
	if err != nil {
		logger.Error(`get namespace by name`, zap.Error(err))

		if errors.Is(err, repository.ErrFindNamespace) {
			return nil, errors.Join(taxonomy.ErrNamespaceNotFound, fmt.Errorf(`%q: %w`, name, err))
		}

		return nil, fmt.Errorf(`unknown error %w`, err)
	}

	if len(nss) != 1 {
		return nil, fmt.Errorf(`%w, got %d results`, taxonomy.ErrNamespaceNotFound, len(nss))
	}

	return nss[0], nil
}
