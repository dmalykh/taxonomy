package loader

import (
	"context"
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
		return nil, err
	}

	// Init zap logger
	var logger *zap.Logger
	switch verbose {
	case true:
		logger, err = zap.NewDevelopment()
		break
	default:
		logger, err = zap.NewProduction()
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
	var transaction = entgo.Transactioner(client)

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
