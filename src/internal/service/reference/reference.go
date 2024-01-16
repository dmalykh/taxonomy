package reference

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
	NamespaceService    taxonomy.Namespace
	referenceRepository repository.Reference
	Logger              *zap.Logger
}

type Service struct {
	log                 *zap.Logger
	namespaceService    taxonomy.Namespace
	referenceRepository repository.Reference
	termService         taxonomy.Term
}

func New(config *Config) taxonomy.Reference {
	return &Service{}
}

func (r *Service) Create(ctx context.Context, termID uint64, namespace string, entitiesID ...model.EntityID) error {
	logger := r.log.With(zap.String(`method`, `SetReference`),
		zap.Uint64(`termID`, termID), zap.String("namespace", namespace), zap.Any(`entitiesID`, entitiesID))

	ns, err := r.namespaceService.GetByName(ctx, namespace)
	if err != nil {
		return fmt.Errorf(`%w %s`, taxonomy.ErrTermNamespaceNotFound, err.Error())
	}

	// Check term exists
	if _, err := r.termService.GetByID(ctx, termID); err != nil {
		logger.Error(`get term by id`, zap.Error(err), zap.Uint64(`termID`, termID))

		if errors.Is(err, repository.ErrFindTerm) {
			return fmt.Errorf(`%w %d`, taxonomy.ErrTermNotFound, termID)
		}

		return fmt.Errorf(`unknown term error %w`, err)
	}

	// Prepare references without duplicates
	var (
		references = make([]*repository.ReferenceModel, 0, len(entitiesID))
		seen       = make(map[model.EntityID]struct{})
	)

	for _, entityID := range entitiesID {
		if _, exists := seen[entityID]; exists {
			continue
		}

		references = append(references, &repository.ReferenceModel{
			TermID:      termID,
			NamespaceID: ns.ID,
			EntityID:    entityID,
		})
		seen[entityID] = struct{}{}
	}

	// Delete exists and insert prepared
	if err := r.referenceRepository.Delete(ctx, []uint64{termID}, []uint64{ns.ID}, entitiesID); err != nil {
		return fmt.Errorf(`can't remove reference %w: %s`, taxonomy.ErrTermReferenceNotRemoved, err.Error())
	}

	if err := r.referenceRepository.Create(ctx, references...); err != nil {
		return fmt.Errorf(`can't create reference %w: %s`, taxonomy.ErrTermReferenceNotCreated, err.Error())
	}

	return nil
}

func (r *Service) Delete(ctx context.Context, termID uint64, namespace string, entitiesID ...model.EntityID) error {
	logger := r.log.With(zap.String(`method`, `UnsetReference`), zap.Uint64(`termID`, termID),
		zap.String("namespace", namespace), zap.Any(`entitiesID`, entitiesID))

	ns, err := r.namespaceService.GetByName(ctx, namespace)
	if err != nil {
		return fmt.Errorf(`%w %s`, taxonomy.ErrTermNamespaceNotFound, err.Error())
	}

	// Check term exists
	if _, err := r.termService.GetByID(ctx, termID); err != nil {
		logger.Error(`get term by id`, zap.Error(err), zap.Uint64(`termID`, termID))

		if errors.Is(err, repository.ErrFindTerm) {
			return fmt.Errorf(`%w %d`, taxonomy.ErrTermNotFound, termID)
		}

		return fmt.Errorf(`unknown term error %w`, err)
	}
	// Remove references
	if err := r.referenceRepository.Delete(ctx, []uint64{termID}, []uint64{ns.ID}, entitiesID); err != nil {
		return fmt.Errorf(`can't remove reference %w: %s`, taxonomy.ErrTermReferenceNotRemoved, err.Error())
	}

	return nil
}

func (t *Service) GetTermsByEntities(ctx context.Context, namespace string, entities ...model.EntityID) ([]model.Term, error) {
	logger := t.log.With(zap.String(`method`, `GetTermsByEntities`),
		zap.String(`namespace`, namespace), zap.Any(`entities`, entities))

	ns, err := t.namespaceService.GetByName(ctx, namespace)
	logger.Debug(`got namespace`, zap.Any(`namespace`, namespace), zap.Error(err))

	if err != nil {
		return nil, taxonomy.ErrTermNamespaceNotFound
	}

	references, err := t.referenceRepository.Get(ctx, &repository.ReferenceFilter{
		NamespaceID: []uint64{ns.ID},
		EntityID:    entities,
	})
	logger.Debug(`got references`, zap.Any(`references`, references), zap.Error(err))

	if err != nil {
		return nil, fmt.Errorf(`unknown error %w`, err)
	}

	terms := make([]model.Term, 0, len(references))

	for _, reference := range references {
		// Change for one request if there is a lot of reference will be found
		term, err := t.termService.GetByID(ctx, reference.TermID)
		logger.Debug(`got term`, zap.Uint64(`id`, reference.TermID), zap.Any(`term`, term), zap.Error(err))

		if err != nil {
			logger.DPanic(`unknown term in reference`, zap.Error(err))

			continue
		}

		terms = append(terms, term)
	}

	return terms, nil
}

func (t *Service) Get(ctx context.Context, filter *model.ReferenceFilter) ([]*model.Reference, error) {
	logger := t.log.With(zap.String(`method`, `GetReferences`), zap.Any(`filter`, filter))

	// Get ids of references
	namespacesID := make([]uint64, 0, len(filter.Namespace))

	for _, ns := range filter.Namespace {
		namespace, err := t.namespaceService.GetByName(ctx, ns)
		if err != nil {
			return nil, taxonomy.ErrTermNamespaceNotFound
		}

		namespacesID = append(namespacesID, namespace.ID)
	}

	// Get references
	references, err := t.referenceRepository.Get(ctx, &repository.ReferenceFilter{
		TermID:      filter.TermID,
		EntityID:    filter.EntityID,
		NamespaceID: namespacesID,
		AfterID:     filter.AfterID,
		Limit:       filter.Limit,
	})
	logger.Debug(`got references`, zap.Any(`references`, references), zap.Error(err))

	if err != nil {
		return nil, fmt.Errorf(`unknown error %w`, err)
	}

	return func() []*model.Reference {
		return nil
	}(), nil
}
