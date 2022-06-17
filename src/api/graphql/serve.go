package graphql

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/dmalykh/tagservice/api/graphql/generated"
	"github.com/dmalykh/tagservice/api/graphql/service"
	"github.com/dmalykh/tagservice/tagservice"
	"log"
	"net/http"
)

type Config struct {
	Port             string
	TagService       tagservice.Tag
	CategoryService  tagservice.Category
	NamespaceService tagservice.Namespace
	Verbose          bool
}

func Serve(config *Config) error {

	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: service.NewResolver(config.TagService, config.CategoryService, config.NamespaceService),
			},
		),
	)
	if config.Verbose {
		srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
			oc := graphql.GetOperationContext(ctx)
			log.Println(oc.RawQuery)
			return next(ctx)
		})
	}

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to :%s for GraphQL playground", config.Port)
	return http.ListenAndServe(":"+config.Port, nil)
}
