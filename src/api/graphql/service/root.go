//nolint:ireturn
package service

import (
	"github.com/dmalykh/tagservice/api/graphql/generated"
	"github.com/dmalykh/tagservice/tagservice"
)

func NewResolver(tagService tagservice.Tag, categoryService tagservice.Category, namespaceService tagservice.Namespace) generated.ResolverRoot { //nolint:lll
	return &Root{
		queryResolver: &Query{
			tagService:      tagService,
			categoryService: categoryService,
		},
		mutationResolver: &Mutation{
			tagService:      tagService,
			categoryService: categoryService,
		},
		entityResolver: &Entity{
			tagService:      tagService,
			categoryService: categoryService,
		},
		categoryResolver: &Category{
			tagService:      tagService,
			categoryService: categoryService,
		},
		tagResolver: &Tag{
			tagService:       tagService,
			categoryService:  categoryService,
			namespaceService: namespaceService,
		},
	}
}

//goland:noinspection GoUnnecessarilyExportedIdentifiers
type Root struct {
	queryResolver    generated.QueryResolver
	mutationResolver generated.MutationResolver
	entityResolver   generated.EntityResolver
	categoryResolver generated.CategoryResolver
	tagResolver      generated.TagResolver
}

func (r *Root) Category() generated.CategoryResolver {
	return r.categoryResolver
}

func (r *Root) Entity() generated.EntityResolver {
	return r.entityResolver
}

func (r *Root) Mutation() generated.MutationResolver {
	return r.mutationResolver
}

func (r *Root) Query() generated.QueryResolver {
	return r.queryResolver
}

func (r *Root) Tag() generated.TagResolver {
	return r.tagResolver
}
