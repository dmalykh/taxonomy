package graphql

import (
	"context"
	"fmt"
	"github.com/dmalykh/taxonomy/taxonomy"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/dmalykh/taxonomy/api/graphql/generated"
	"github.com/dmalykh/taxonomy/api/graphql/service"
)

type Config struct {
	Port              string
	TermService       taxonomy.Term
	VocabularyService taxonomy.Vocabulary
	NamespaceService  taxonomy.Namespace
	Verbose           bool
}

func Serve(config *Config) error {
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: service.NewResolver(config.TermService, config.VocabularyService, config.NamespaceService),
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

	http.Handle("/", playground.Handler("Taxonomy GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to :%s for GraphQL playground", config.Port)

	return fmt.Errorf(`server error: %w`, http.ListenAndServe(":"+config.Port, nil))
}
