package main

import (
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
}

func Serve(config *Config) {

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: service.NewResolver(config.TagService, config.CategoryService, config.NamespaceService)}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to :%s for GraphQL playground", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
