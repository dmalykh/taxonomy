//nolint:ireturn
package service

import (
	"github.com/dmalykh/taxonomy/api/graphql/generated"
	"github.com/dmalykh/taxonomy/taxonomy"
)

func NewResolver(termService taxonomy.Term, vocabularyService taxonomy.Vocabulary, namespaceService taxonomy.Namespace) generated.ResolverRoot { //nolint:lll
	return &Root{
		queryResolver: &Query{
			termService:       termService,
			vocabularyService: vocabularyService,
		},
		mutationResolver: &Mutation{
			termService:       termService,
			vocabularyService: vocabularyService,
		},
		entityResolver: &Entity{
			termService:       termService,
			vocabularyService: vocabularyService,
		},
		vocabularyResolver: &Vocabulary{
			termService:       termService,
			vocabularyService: vocabularyService,
		},
		termResolver: &Term{
			termService:       termService,
			vocabularyService: vocabularyService,
			namespaceService:  namespaceService,
		},
	}
}

type Root struct {
	queryResolver      generated.QueryResolver
	mutationResolver   generated.MutationResolver
	entityResolver     generated.EntityResolver
	vocabularyResolver generated.VocabularyResolver
	termResolver       generated.TermResolver
}

func (r *Root) Vocabulary() generated.VocabularyResolver {
	return r.vocabularyResolver
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

func (r *Root) Term() generated.TermResolver {
	return r.termResolver
}
