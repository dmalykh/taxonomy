package reference

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/taxonomy/taxonomy"
	"github.com/dmalykh/taxonomy/taxonomy/model"
	"github.com/dmalykh/taxonomy/taxonomy/repository"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type Config struct {
	NamespaceService    taxonomy.Namespace
	ReferenceRepository repository.Reference
	TermService         taxonomy.Term
	Logger              *zap.Logger
}

type Service struct {
	log                 *zap.Logger
	namespaceService    taxonomy.Namespace
	referenceRepository repository.Reference
	termService         taxonomy.Term
}

func New(config *Config) taxonomy.Reference {
	return &Service{
		namespaceService:    config.NamespaceService,
		termService:         config.TermService,
		referenceRepository: config.ReferenceRepository,
		log:                 config.Logger,
	}
}

func (r *Service) Create(ctx context.Context, termID uint64, namespace string, entitiesID ...model.EntityID) error {
	logger := r.log.With(zap.String(`method`, `SetReference`),
		zap.Uint64(`termID`, termID), zap.String("namespace", namespace), zap.Any(`entitiesID`, entitiesID))

	ns, err := r.namespaceService.GetByName(ctx, namespace)
	if err != nil {
		return fmt.Errorf(`namespace %s get error: %w: %w`, namespace, taxonomy.ErrNamespaceNotFound, err)
	}

	// Check term exists
	if _, err := r.termService.GetByID(ctx, termID); err != nil {
		logger.Error(`get term by id`, zap.Error(err), zap.Uint64(`termID`, termID))

		if errors.Is(err, repository.ErrFindTerm) {
			return fmt.Errorf(`%w %d`, taxonomy.ErrTermNotFound, termID)
		}

		return fmt.Errorf(`unknown term %d error %w`, termID, err)
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

	// Upsert prepared
	if err := r.referenceRepository.Set(ctx, references...); err != nil {
		return fmt.Errorf(`can't create reference %w: %w`, taxonomy.ErrReferenceNotCreated, err)
	}

	return nil
}

func (r *Service) Delete(ctx context.Context, termID uint64, namespace string, entitiesID ...model.EntityID) error {
	r.log.With(zap.String(`method`, `Delete`), zap.Uint64(`termID`, termID),
		zap.String("namespace", namespace), zap.Any(`entitiesID`, entitiesID)).Info(`delete reference`)

	ns, err := r.namespaceService.GetByName(ctx, namespace)
	if err != nil {
		return fmt.Errorf(`namespace %s get error: %w: %w`, namespace, taxonomy.ErrNamespaceNotFound, err)
	}

	// Remove references
	if err := r.referenceRepository.Delete(ctx, &repository.ReferenceFilter{
		TermID:      [][]uint64{{termID}},
		NamespaceID: []uint64{ns.ID},
		EntityID:    entitiesID,
	}); err != nil {
		return fmt.Errorf(`can't remove reference %w: %w`, taxonomy.ErrReferenceNotRemoved, err)
	}

	return nil
}

// Todo: Do we need this method?
//func (t *Service) GetTerms(ctx context.Context, namespace string, entities ...model.EntityID) ([]model.Term, error) {
//	logger := t.log.With(zap.String(`method`, `GetTerms`),
//		zap.String(`namespace`, namespace), zap.Any(`entities`, entities))
//
//	ns, err := t.namespaceService.GetByName(ctx, namespace)
//	if err != nil {
//		return nil, taxonomy.ErrNamespaceNotFound
//	}
//
//	references, err := t.referenceRepository.Get(ctx, &repository.ReferenceFilter{
//		NamespaceID: []uint64{ns.ID},
//		EntityID:    entities,
//	})
//	logger.Debug(`got references`, zap.Any(`references`, references), zap.Error(err))
//
//	if err != nil {
//		return nil, fmt.Errorf(`unknown error %w`, err)
//	}
//
//	terms := make([]model.Term, 0, len(references))
//
//	for _, reference := range references {
//		// Change for one request if there is a lot of reference will be found
//		term, err := t.termService.GetByID(ctx, reference.TermID)
//		logger.Debug(`got term`, zap.Uint64(`id`, reference.TermID), zap.Any(`term`, term), zap.Error(err))
//
//		if err != nil {
//			logger.DPanic(`unknown term in reference`, zap.Error(err))
//
//			continue
//		}
//
//		//terms = append(terms, term)
//	}
//
//	return terms, nil
//}

func (t *Service) Get(ctx context.Context, filter *model.ReferenceFilter) ([]*model.Reference, error) {
	logger := t.log.With(zap.String(`method`, `GetReferences`), zap.Any(`filter`, filter))

	// Get ids of references
	namespaces := make(map[uint64]*model.Namespace, len(filter.Namespace))

	for _, ns := range filter.Namespace {
		ns, err := t.namespaceService.GetByName(ctx, ns)
		if err != nil {
			return nil, fmt.Errorf(`%w: %w`, taxonomy.ErrNamespaceNotFound, err)
		}

		namespaces[ns.ID] = ns
	}

	// Get references
	references, err := t.referenceRepository.Get(ctx, &repository.ReferenceFilter{
		TermID:      filter.TermID,
		EntityID:    filter.EntityID,
		NamespaceID: lo.Keys[uint64, *model.Namespace](namespaces),
		AfterID:     filter.AfterID,
		Limit:       filter.Limit,
	})
	logger.Debug(`got references`, zap.Any(`references`, references), zap.Error(err))

	if err != nil {
		return nil, fmt.Errorf(`unknown error %w`, err)
	}

	return func(references []*repository.ReferenceModel) []*model.Reference {
		var models = make([]*model.Reference, 0, len(references))
		for _, ref := range references {
			models = append(models, &model.Reference{
				ID:        ref.ID,
				TermID:    ref.TermID,
				Namespace: namespaces[ref.NamespaceID].Data.Name,
				EntityID:  ref.EntityID,
			})
		}

		return models
	}(references), nil
}
