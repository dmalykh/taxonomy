package term_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/taxonomy/taxonomy"
	"github.com/dmalykh/taxonomy/taxonomy/model"
	"github.com/dmalykh/taxonomy/taxonomy/repository"
	"go.uber.org/zap"
	"io"
	"testing"

	"github.com/dmalykh/taxonomy/internal/service/term"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/stretchr/testify/assert"
)

func TestTermService_Delete(t *testing.T) {
	t.Parallel()

	errunknown := errors.New(`unknown`)

	tests := []struct {
		name        string
		TermService func(ctx context.Context) taxonomy.Term
		err         error
	}{
		{
			name: `not found term error`,
			TermService: func(ctx context.Context) taxonomy.Term {
				termrepo := mock.Mock[repository.Term]()
				mock.When(termrepo.GetByID(mock.Exact[context.Context](ctx), mock.Any[uint64]())).
					ThenReturn(model.Term{}, repository.ErrFindTerm)

				ref := mock.Mock[taxonomy.Reference]()

				return term.New(&term.Config{
					ReferenceService: ref,
					TermRepository:   termrepo,
					Logger:           zap.NewNop(),
				})
			},
			err: taxonomy.ErrTermNotFound,
		},
		{
			name: `unknown err when getting term`,
			TermService: func(ctx context.Context) taxonomy.Term {
				termrepo := mock.Mock[repository.Term]()
				mock.When(termrepo.GetByID(mock.Exact[context.Context](ctx), mock.Any[uint64]())).
					ThenReturn(model.Term{}, errunknown)

				ref := mock.Mock[taxonomy.Reference]()

				return term.New(&term.Config{
					ReferenceService: ref,
					TermRepository:   termrepo,
					Logger:           zap.NewNop(),
				})
			},
			err: errunknown,
		},
		{
			name: `error find reference`,
			TermService: func(ctx context.Context) taxonomy.Term {
				termrepo := mock.Mock[repository.Term]()
				mock.When(termrepo.GetByID(mock.Exact[context.Context](ctx), mock.Any[uint64]())).
					ThenReturn(model.Term{ID: 33}, nil)

				ref := mock.Mock[taxonomy.Reference]()
				mock.When(ref.Get(mock.Exact[context.Context](ctx), mock.Any[*model.ReferenceFilter]())).
					ThenReturn(nil, errunknown)

				return term.New(&term.Config{
					ReferenceService: ref,
					TermRepository:   termrepo,
					Logger:           zap.NewNop(),
				})
			},
			err: errunknown,
		},
		{
			name: `error non-empty reference`,
			TermService: func(ctx context.Context) taxonomy.Term {
				termrepo := mock.Mock[repository.Term]()
				mock.When(termrepo.GetByID(mock.Exact[context.Context](ctx), mock.Any[uint64]())).
					ThenReturn(model.Term{}, nil)

				ref := mock.Mock[taxonomy.Reference]()
				mock.When(ref.Get(mock.Exact[context.Context](ctx), mock.Any[*model.ReferenceFilter]())).
					ThenReturn([]model.Reference{{ID: 9}}, nil)

				return term.New(&term.Config{
					ReferenceService: ref,
					TermRepository:   termrepo,
					Logger:           zap.NewNop(),
				})
			},
			err: taxonomy.ErrTermReferenceExists,
		},
		{
			name: `error remove term`,
			TermService: func(ctx context.Context) taxonomy.Term {
				termrepo := mock.Mock[repository.Term]()
				mock.When(termrepo.GetByID(mock.Exact[context.Context](ctx), mock.Any[uint64]())).
					ThenReturn(model.Term{}, nil)

				ref := mock.Mock[taxonomy.Reference]()
				mock.When(ref.Get(mock.Exact[context.Context](ctx), mock.Any[*model.ReferenceFilter]())).
					ThenReturn(nil, nil)

				mock.When(termrepo.DeleteByID(mock.Exact[context.Context](ctx), mock.Any[uint64]())).
					ThenReturn(io.EOF)

				return term.New(&term.Config{
					ReferenceService: ref,
					TermRepository:   termrepo,
					Logger:           zap.NewNop(),
				})
			},
			err: io.EOF,
		},
		{
			name: `no error`,
			TermService: func(ctx context.Context) taxonomy.Term {
				termrepo := mock.Mock[repository.Term]()
				mock.When(termrepo.GetByID(mock.Exact[context.Context](ctx), mock.Any[uint64]())).
					ThenReturn(model.Term{}, nil)

				ref := mock.Mock[taxonomy.Reference]()
				mock.When(ref.Get(mock.Exact[context.Context](ctx), mock.Any[*model.ReferenceFilter]())).
					ThenReturn(nil, nil)

				mock.When(termrepo.DeleteByID(mock.Exact[context.Context](ctx), mock.Any[uint64]())).
					ThenReturn(nil)

				return term.New(&term.Config{
					ReferenceService: ref,
					TermRepository:   termrepo,
					Logger:           zap.NewNop(),
				})
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mock.SetUp(t)

			var ctx = context.Background()

			// Initialize service
			s := tt.TermService(ctx)

			err := s.Delete(ctx, 88)
			if tt.err == nil {
				assert.NoError(t, err)

				return
			}
			assert.ErrorIs(t, err, tt.err)
		})
	}
}

func TestTermService_Update(t *testing.T) {

	var defaultTermData = model.TermData{
		Name:         `utrecht`,
		Title:        `Utrecht`,
		Description:  `Utrecht is a city in the Netherlands, known for its medieval center`,
		VocabularyID: []uint64{4, 8, 15},
	}

	tests := []struct {
		name        string
		TermService func(ctx context.Context) taxonomy.Term
		update      *model.TermData
		want        model.Term
		err         error
	}{
		{
			name: `updated fine`,
			TermService: func(ctx context.Context) taxonomy.Term {
				termrepo := mock.Mock[repository.Term]()
				mock.When(termrepo.GetByID(mock.Exact[context.Context](ctx), mock.Any[uint64]())).
					ThenReturn(model.Term{}, nil)
				mock.When(termrepo.Update(mock.Exact[context.Context](ctx), mock.Any[uint64](), mock.Any[*model.TermData]())).
					ThenReturn(model.Term{ID: 22, Data: defaultTermData}, nil)

				voc := mock.Mock[taxonomy.Vocabulary]()
				mock.When(voc.GetByID(mock.Exact[context.Context](ctx), mock.Any[uint64]())).
					ThenReturn(model.Vocabulary{}, nil)

				return term.New(&term.Config{
					TermRepository:    termrepo,
					VocabularyService: voc,
					Logger:            zap.NewNop(),
				})
			},
			update: &defaultTermData,
			want:   model.Term{ID: 22, Data: defaultTermData},
			err:    nil,
		},
		{
			name: `updated fine without data`,
			TermService: func(ctx context.Context) taxonomy.Term {
				termrepo := mock.Mock[repository.Term]()
				mock.When(termrepo.GetByID(mock.Exact[context.Context](ctx), mock.Any[uint64]())).
					ThenReturn(model.Term{}, nil)
				mock.When(termrepo.Update(mock.Exact[context.Context](ctx), mock.Any[uint64](), mock.Any[*model.TermData]())).
					ThenReturn(model.Term{ID: 22, Data: defaultTermData}, nil)

				return term.New(&term.Config{
					TermRepository: termrepo,
					Logger:         zap.NewNop(),
				})
			},
			update: &model.TermData{},
			want:   model.Term{ID: 22, Data: defaultTermData},
			err:    nil,
		},
		{
			name: `error update`,
			TermService: func(ctx context.Context) taxonomy.Term {
				termrepo := mock.Mock[repository.Term]()
				mock.When(termrepo.GetByID(mock.Exact[context.Context](ctx), mock.Any[uint64]())).
					ThenReturn(model.Term{}, nil)
				mock.When(termrepo.Update(mock.Exact[context.Context](ctx), mock.Any[uint64](), mock.Any[*model.TermData]())).
					ThenReturn(model.Term{}, io.EOF)

				return term.New(&term.Config{
					TermRepository: termrepo,
					Logger:         zap.NewNop(),
				})
			},
			update: &model.TermData{},
			want:   model.Term{},
			err:    taxonomy.ErrTermNotUpdated,
		},
		{
			name: `error term doesn't exists`,
			TermService: func(ctx context.Context) taxonomy.Term {
				termrepo := mock.Mock[repository.Term]()
				mock.When(termrepo.GetByID(mock.Exact[context.Context](ctx), mock.Any[uint64]())).
					ThenReturn(model.Term{}, fmt.Errorf(`%w`, repository.ErrFindTerm))

				return term.New(&term.Config{
					TermRepository: termrepo,
					Logger:         zap.NewNop(),
				})
			},
			update: &defaultTermData,
			want:   model.Term{},
			err:    taxonomy.ErrTermNotFound,
		},
		{
			name: `error get term`,
			TermService: func(ctx context.Context) taxonomy.Term {
				termrepo := mock.Mock[repository.Term]()
				mock.When(termrepo.GetByID(mock.Exact[context.Context](ctx), mock.Any[uint64]())).
					ThenReturn(model.Term{}, io.EOF)

				return term.New(&term.Config{
					TermRepository: termrepo,
					Logger:         zap.NewNop(),
				})
			},
			update: &defaultTermData,
			want:   model.Term{},
			err:    io.EOF,
		},
		{
			name: `error check second id in vocabulary`,
			TermService: func(ctx context.Context) taxonomy.Term {
				termrepo := mock.Mock[repository.Term]()
				mock.When(termrepo.GetByID(mock.Exact[context.Context](ctx), mock.Any[uint64]())).
					ThenReturn(model.Term{}, nil)

				voc := mock.Mock[taxonomy.Vocabulary]()
				mock.When(voc.GetByID(mock.Exact[context.Context](ctx), mock.NotEqual[uint64](99))).
					ThenReturn(model.Vocabulary{}, nil)
				mock.When(voc.GetByID(mock.Exact[context.Context](ctx), mock.Equal[uint64](99))).
					ThenReturn(model.Vocabulary{}, repository.ErrFindVocabulary)

				return term.New(&term.Config{
					TermRepository:    termrepo,
					VocabularyService: voc,
					Logger:            zap.NewNop(),
				})
			},
			update: &model.TermData{
				VocabularyID: []uint64{4, 8, 99, 15},
			},
			want: model.Term{},
			err:  taxonomy.ErrVocabularyNotFound,
		},
		{
			name: `error during vocabularies check`,
			TermService: func(ctx context.Context) taxonomy.Term {
				termrepo := mock.Mock[repository.Term]()
				mock.When(termrepo.GetByID(mock.Exact[context.Context](ctx), mock.Any[uint64]())).
					ThenReturn(model.Term{}, nil)

				voc := mock.Mock[taxonomy.Vocabulary]()
				mock.When(voc.GetByID(mock.Exact[context.Context](ctx), mock.NotEqual[uint64](99))).
					ThenReturn(model.Vocabulary{}, nil)
				mock.When(voc.GetByID(mock.Exact[context.Context](ctx), mock.Equal[uint64](99))).
					ThenReturn(model.Vocabulary{}, io.EOF)

				return term.New(&term.Config{
					TermRepository:    termrepo,
					VocabularyService: voc,
					Logger:            zap.NewNop(),
				})
			},
			update: &model.TermData{
				VocabularyID: []uint64{4, 8, 99, 15},
			},
			want: model.Term{},
			err:  io.EOF,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.SetUp(t)
			var ctx = context.Background()
			s := tt.TermService(ctx)

			got, err := s.Update(ctx, 99, tt.update)
			assert.Equal(t, tt.want, got)
			if tt.err == nil {
				assert.NoError(t, err)
				return
			}
			assert.ErrorIs(t, err, tt.err)
		})
	}
}

//
//func TestTermService_GetReferences(t *testing.T) {
//	t.Parallel()
//
//	tests := []struct {
//		name                      string
//		termGroups                [][]uint
//		ReferenceGetReturns       func() ([]model.Reference, error)
//		NamespaceGetByNameReturns func() (model.NamespaceID, error)
//		want                      func(t assert.TestingT, references []model.Reference)
//		wantErr                   assert.ErrorAssertionFunc
//	}{
//		{
//			name:       `no references`,
//			termGroups: [][]uint{{34}, {92, 96}},
//			NamespaceGetByNameReturns: func() (model.NamespaceID, error) {
//				return model.NamespaceID{ID: 11}, nil
//			},
//			ReferenceGetReturns: func() ([]model.Reference, error) {
//				return nil, nil
//			},
//			want: func(t assert.TestingT, references []model.Reference) {
//				assert.Len(t, references, 0)
//			},
//			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
//				return assert.NoError(t, err)
//			},
//		},
//		{
//			name:       `namespace not found`,
//			termGroups: [][]uint{{34}, {92, 96}},
//			NamespaceGetByNameReturns: func() (model.NamespaceID, error) {
//				return model.NamespaceID{ID: 11}, errors.New(``)
//			},
//			ReferenceGetReturns: func() ([]model.Reference, error) {
//				return nil, nil
//			},
//			want: func(t assert.TestingT, references []model.Reference) {
//				assert.Len(t, references, 0)
//			},
//			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
//				return assert.ErrorIs(t, err, taxonomy.ErrTermNamespaceNotFound)
//			},
//		},
//	}
//
//	// Create categories "RAM", "CPU", "Display size"
//	// Should receive all laptops that has: "RAM" (512 or 1024) and "CPU" (2.8 or 3.2) and "Display size"
//	// (between 13 and 15)
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			t.Parallel()
//			referencerepo := mockrepository.NewReference(t)
//			referencerepo.On(`Get`, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
//				Return(tt.ReferenceGetReturns()).Maybe()
//
//			namespaceservice := mockservice.NewNamespace(t)
//			namespaceservice.On(`GetByName`, mock.Anything, mock.Anything).Return(tt.NamespaceGetByNameReturns())
//
//			termService := term.New(&term.Config{
//				NamespaceService:    namespaceservice,
//				ReferenceRepository: referencerepo,
//				Logger:              zap.NewNop(),
//			})
//
//			references, err := termService.GetReferences(context.TODO(), &model.EntityFilter{
//				TermID:    tt.termGroups,
//				NamespaceID: []string{`any`},
//			})
//			tt.wantErr(t, err)
//			tt.want(t, references)
//		})
//	}
//}
