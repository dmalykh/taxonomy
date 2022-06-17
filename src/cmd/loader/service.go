package loader

import (
	"context"
	"fmt"

	"github.com/dmalykh/tagservice/repository/entgo"
	"github.com/dmalykh/tagservice/repository/entgo/repository"
	"github.com/dmalykh/tagservice/tagservice"
	"github.com/dmalykh/tagservice/tagservice/service/category"
	"github.com/dmalykh/tagservice/tagservice/service/namespace"
	"github.com/dmalykh/tagservice/tagservice/service/tag"
	"go.uber.org/zap"
)

type Service struct {
	Namespace tagservice.Namespace
	Tag       tagservice.Tag
	Category  tagservice.Category
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

	transaction := entgo.Transactioner(client)

	service.Namespace = namespace.New(&namespace.Config{
		Transaction:         transaction,
		NamespaceRepository: repository.NewNamespace(client.Namespace),
		RelationRepository:  repository.NewRelation(client.Relation),
		Logger:              logger,
	})

	service.Tag = tag.New(&tag.Config{
		Transaction:        transaction,
		TagRepository:      repository.NewTag(client.Tag),
		CategoryRepository: repository.NewCategory(client.Category),
		RelationRepository: repository.NewRelation(client.Relation),
		NamespaceService:   service.Namespace,
		Logger:             logger,
	})

	service.Category = category.New(&category.Config{
		CategoryRepository: repository.NewCategory(client.Category),
		TagService:         service.Tag,
		Logger:             logger,
	})

	return &service, nil
}
