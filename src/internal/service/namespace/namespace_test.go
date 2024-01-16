package namespace_test

import (
	"context"
	"errors"
	"github.com/dmalykh/taxonomy/internal/service/namespace"
	"github.com/dmalykh/taxonomy/taxonomy"
	"github.com/dmalykh/taxonomy/taxonomy/model"
	"github.com/dmalykh/taxonomy/taxonomy/repository"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestNamespaceService_Create(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run(`create with error`, func(t *testing.T) {
		var ctx = context.Background()
		mock.SetUp(t)
		repo := mock.Mock[repository.Namespace]()
		s := namespace.New(&namespace.Config{
			Logger:              logger,
			NamespaceRepository: repo,
		})

		mock.When(repo.Create(mock.Exact[context.Context](ctx), mock.Any[*model.NamespaceData]())).
			ThenAnswer(func(args []any) []any {
				assert.Equal(t, `ambar`, args[1].(*model.NamespaceData).Name)
				return []any{&model.Namespace{}, errors.New(``)}
			})

		_, err := s.Create(ctx, `ambar`)
		assert.Error(t, err)
		mock.Verify(repo, mock.Once())
	})

	t.Run(`create without error`, func(t *testing.T) {
		var ctx = context.Background()
		mock.SetUp(t)
		repo := mock.Mock[repository.Namespace]()
		s := namespace.New(&namespace.Config{
			Logger:              logger,
			NamespaceRepository: repo,
		})
		mock.When(repo.Create(mock.Exact[context.Context](ctx), mock.Any[*model.NamespaceData]())).
			ThenAnswer(func(args []any) []any {
				assert.Equal(t, `zamok`, args[1].(*model.NamespaceData).Name)
				return []any{&model.Namespace{}, nil}
			})

		_, err := s.Create(ctx, `zamok`)
		assert.NoError(t, err)
		mock.Verify(repo, mock.Once())
	})
}

func TestNamespaceService_Delete(t *testing.T) {
	t.Parallel()

	errunknown := errors.New(`unknown`)

	tests := []struct {
		name    string
		id      uint64
		Service func(ctx context.Context) taxonomy.Namespace
		err     error
	}{
		{
			name: `not found namespace error`,
			id:   88,
			Service: func(ctx context.Context) taxonomy.Namespace {
				var namespacerepo = mock.Mock[repository.Namespace]()
				mock.When(namespacerepo.Get(mock.Exact[context.Context](ctx), mock.Any[*repository.NamespaceFilter]())).
					ThenAnswer(func(args []any) []any {
						assert.Equal(t, uint64(88), args[1].(*repository.NamespaceFilter).ID[0])
						return []any{nil, repository.ErrFindNamespace}
					})

				return namespace.New(&namespace.Config{
					Logger:              zap.NewNop(),
					NamespaceRepository: namespacerepo,
				})
			},
			err: taxonomy.ErrNamespaceNotFound,
		},
		{
			name: `unknown err when getting namespace`,
			id:   99,
			Service: func(ctx context.Context) taxonomy.Namespace {
				var namespacerepo = mock.Mock[repository.Namespace]()
				mock.When(namespacerepo.Get(mock.Exact[context.Context](ctx), mock.Any[*repository.NamespaceFilter]())).
					ThenAnswer(func(args []any) []any {
						assert.Equal(t, uint64(99), args[1].(*repository.NamespaceFilter).ID[0])
						return []any{nil, errunknown}
					})

				return namespace.New(&namespace.Config{
					Logger:              zap.NewNop(),
					NamespaceRepository: namespacerepo,
				})
			},
			err: errunknown,
		},
		{
			name: `error find reference`,
			id:   100,
			Service: func(ctx context.Context) taxonomy.Namespace {
				var namespacerepo = mock.Mock[repository.Namespace]()
				mock.When(namespacerepo.Get(mock.Exact[context.Context](ctx), mock.Any[*repository.NamespaceFilter]())).
					ThenAnswer(func(args []any) []any {
						assert.Equal(t, uint64(100), args[1].(*repository.NamespaceFilter).ID[0])
						return []any{[]*model.Namespace{
							{ID: 100},
						}, nil}
					})

				var ref = mock.Mock[taxonomy.Reference]()
				mock.When(ref.Get(mock.Exact[context.Context](ctx), mock.Any[*model.ReferenceFilter]())).
					ThenReturn(nil, errunknown)

				return namespace.New(&namespace.Config{
					Logger:              zap.NewNop(),
					ReferenceService:    ref,
					NamespaceRepository: namespacerepo,
				})
			},
			err: errunknown,
		},
		{
			name: `error non-empty reference`,
			id:   444,
			Service: func(ctx context.Context) taxonomy.Namespace {
				var namespacerepo = mock.Mock[repository.Namespace]()
				mock.When(namespacerepo.Get(mock.Exact[context.Context](ctx), mock.Any[*repository.NamespaceFilter]())).
					ThenAnswer(func(args []any) []any {
						assert.Equal(t, uint64(444), args[1].(*repository.NamespaceFilter).ID[0])
						return []any{[]*model.Namespace{
							{ID: 444},
						}, nil}
					})

				var ref = mock.Mock[taxonomy.Reference]()
				mock.When(ref.Get(mock.Exact[context.Context](ctx), mock.Any[*model.ReferenceFilter]())).
					ThenReturn([]*model.Reference{
						{ID: 444},
					}, nil)

				return namespace.New(&namespace.Config{
					Logger:              zap.NewNop(),
					ReferenceService:    ref,
					NamespaceRepository: namespacerepo,
				})
			},
			err: taxonomy.ErrTermReferenceExists,
		},
		{
			name: `error delete namespace`,
			id:   555,
			Service: func(ctx context.Context) taxonomy.Namespace {
				var namespacerepo = mock.Mock[repository.Namespace]()
				mock.When(namespacerepo.Get(mock.Exact[context.Context](ctx), mock.Any[*repository.NamespaceFilter]())).
					ThenAnswer(func(args []any) []any {
						assert.Equal(t, uint64(555), args[1].(*repository.NamespaceFilter).ID[0])
						return []any{[]*model.Namespace{
							{ID: 555},
						}, nil}
					})
				mock.When(namespacerepo.Delete(mock.Exact[context.Context](ctx), mock.Any[*repository.NamespaceFilter]())).
					ThenAnswer(func(args []any) []any {
						assert.Equal(t, uint64(555), args[1].(*repository.NamespaceFilter).ID[0])
						return []any{errunknown}
					})

				var ref = mock.Mock[taxonomy.Reference]()
				mock.When(ref.Get(mock.Exact[context.Context](ctx), mock.Any[*model.ReferenceFilter]())).
					ThenReturn(nil, nil)

				return namespace.New(&namespace.Config{
					Logger:              zap.NewNop(),
					ReferenceService:    ref,
					NamespaceRepository: namespacerepo,
				})
			},
			err: errunknown,
		},
		{
			name: `no error`,
			id:   666,
			Service: func(ctx context.Context) taxonomy.Namespace {
				var namespacerepo = mock.Mock[repository.Namespace]()
				mock.When(namespacerepo.Get(mock.Exact[context.Context](ctx), mock.Any[*repository.NamespaceFilter]())).
					ThenAnswer(func(args []any) []any {
						assert.Equal(t, uint64(666), args[1].(*repository.NamespaceFilter).ID[0])
						return []any{[]*model.Namespace{
							{ID: 666},
						}, nil}
					})
				mock.When(namespacerepo.Delete(mock.Exact[context.Context](ctx), mock.Any[*repository.NamespaceFilter]())).
					ThenAnswer(func(args []any) []any {
						assert.Equal(t, uint64(666), args[1].(*repository.NamespaceFilter).ID[0])
						return []any{nil}
					})

				var ref = mock.Mock[taxonomy.Reference]()
				mock.When(ref.Get(mock.Exact[context.Context](ctx), mock.Any[*model.ReferenceFilter]())).
					ThenReturn(nil, nil)

				return namespace.New(&namespace.Config{
					Logger:              zap.NewNop(),
					ReferenceService:    ref,
					NamespaceRepository: namespacerepo,
				})
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mock.SetUp(t)
			var ctx = context.Background()
			s := tt.Service(ctx)

			err := s.Delete(ctx, tt.id)
			if tt.err == nil {
				assert.NoError(t, err)
				return
			}
			assert.ErrorIs(t, err, tt.err)
		})
	}
}
