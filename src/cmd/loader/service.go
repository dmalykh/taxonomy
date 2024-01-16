package loader

import (
	"context"
	"fmt"
	"github.com/dmalykh/taxonomy/internal/repository/entgo"
	repository2 "github.com/dmalykh/taxonomy/internal/repository/entgo/repository"
	"github.com/dmalykh/taxonomy/taxonomy"

	"github.com/dmalykh/taxonomy/internal/service/namespace"
	"github.com/dmalykh/taxonomy/internal/service/term"
	"github.com/dmalykh/taxonomy/internal/service/vocabulary"
	"go.uber.org/zap"
)

type Service struct {
	Namespace  taxonomy.Namespace
	Term       taxonomy.Term
	Vocabulary taxonomy.Vocabulary
}

func Load(ctx context.Context, dsn string, verbose bool) (*Service, error) {
	client, err := entgo.Connect(ctx, dsn, verbose)
	if err != nil {
		return nil, fmt.Errorf(`error connect to database: %w`, err)
	}

	// Init zap logger
	logger, err := func() (*zap.Logger, error) {
		if verbose == true {
			return zap.NewDevelopment() //nolint:wrapcheck
		}

		return zap.NewProduction() //nolint:wrapcheck
	}()
	if err != nil {
		return nil, fmt.Errorf(`logger error: %w`, err)
	}

	// Graceful shutdown
	go func() {
		<-ctx.Done()
		func() {
			defer func() {
				if r := recover(); r != nil {
					logger.DPanic(`panic`, zap.Any(`recover`, r))
				}
			}()

			err := client.Close()
			if err != nil {
				panic(err)
			}
		}()
		func() {
			err := logger.Sync()
			if err != nil {
				logger.DPanic(`panic`, zap.Error(err))
			}
		}()
	}()

	// Construct service
	var service Service

	service.Namespace = namespace.New(&namespace.Config{
		Transaction:         transaction,
		NamespaceRepository: repository2.NewNamespace(client.Namespace),
		ReferenceRepository: repository2.NewReference(client.Reference),
		Logger:              logger,
	})

	service.Term = term.New(&term.Config{
		Transaction:         transaction,
		TermRepository:      repository2.NewTerm(client.Term),
		VocabularyService:   repository2.NewVocabulary(client.Vocabulary),
		ReferenceRepository: repository2.NewReference(client.Reference),
		NamespaceService:    service.Namespace,
		Logger:              logger,
	})

	service.Vocabulary = vocabulary.New(&vocabulary.Config{
		VocabularyRepository: repository2.NewVocabulary(client.Vocabulary),
		TermService:          service.Term,
		Logger:               logger,
	})

	return &service, nil
}
