package reference_test

import (
	"context"
	"github.com/dmalykh/taxonomy/internal/service/reference"
	"github.com/dmalykh/taxonomy/taxonomy"
	"github.com/dmalykh/taxonomy/taxonomy/model"
	"github.com/dmalykh/taxonomy/taxonomy/repository"
	"github.com/jaswdr/faker"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io"
	"testing"
)

func TestService_Create(t *testing.T) {

	type args struct {
		termID     uint64
		namespace  string
		entitiesID []model.EntityID
	}
	tests := []struct {
		name   string
		args   args
		check  func(t *testing.T, err error)
		config func() *reference.Config
	}{
		{
			name: `namespace not found`,
			args: args{
				termID:    6,
				namespace: `kleo`,
			},
			config: func() *reference.Config {
				ns := mock.Mock[taxonomy.Namespace]()
				mock.When(ns.GetByName(mock.Any[context.Context](), mock.Exact[string](`kleo`))).
					ThenReturn(nil, io.EOF)

				return &reference.Config{
					NamespaceService: ns,
				}
			},
			check: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, taxonomy.ErrNamespaceNotFound)
				assert.ErrorIs(t, err, io.EOF)
				assert.ErrorContains(t, err, `kleo`)
			},
		},
		{
			name: `term not found`,
			args: args{
				termID:    138,
				namespace: `kleo`,
			},
			config: func() *reference.Config {
				ns := mock.Mock[taxonomy.Namespace]()
				mock.When(ns.GetByName(mock.Any[context.Context](), mock.Exact[string](`kleo`))).
					ThenReturn(&model.Namespace{ID: 66}, nil)

				trm := mock.Mock[taxonomy.Term]()
				mock.When(trm.GetByID(mock.Any[context.Context](), mock.Exact[uint64](138))).
					ThenReturn(nil, repository.ErrFindTerm)

				return &reference.Config{
					NamespaceService: ns,
					TermService:      trm,
				}
			},
			check: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, taxonomy.ErrTermNotFound)
				assert.ErrorContains(t, err, `138`)
			},
		},
		{
			name: `get term error`,
			args: args{
				termID:    138,
				namespace: `kleo`,
			},
			config: func() *reference.Config {
				ns := mock.Mock[taxonomy.Namespace]()
				mock.When(ns.GetByName(mock.Any[context.Context](), mock.Exact[string](`kleo`))).
					ThenReturn(&model.Namespace{ID: 66}, nil)

				trm := mock.Mock[taxonomy.Term]()
				mock.When(trm.GetByID(mock.Any[context.Context](), mock.Exact[uint64](138))).
					ThenReturn(nil, io.EOF)

				return &reference.Config{
					NamespaceService: ns,
					TermService:      trm,
				}
			},
			check: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, io.EOF)
				assert.ErrorContains(t, err, `138`)
			},
		},
		{
			name: `set reference error`,
			args: args{
				termID:     138,
				namespace:  `kleo`,
				entitiesID: []model.EntityID{`smart`, `beautiful`, `pretty`, `clever`, `strong`},
			},
			config: func() *reference.Config {
				ns := mock.Mock[taxonomy.Namespace]()
				mock.When(ns.GetByName(mock.Any[context.Context](), mock.Exact[string](`kleo`))).
					ThenReturn(&model.Namespace{ID: 66}, nil)

				trm := mock.Mock[taxonomy.Term]()
				mock.When(trm.GetByID(mock.Any[context.Context](), mock.Exact[uint64](138))).
					ThenReturn(&model.Term{}, nil)

				ref := mock.Mock[repository.Reference]()
				mock.When(ref.Set(mock.Any[context.Context](), mock.Any[[]*repository.ReferenceModel]()...)).
					ThenReturn(repository.ErrCreateReference)

				return &reference.Config{
					NamespaceService:    ns,
					TermService:         trm,
					ReferenceRepository: ref,
				}
			},
			check: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, repository.ErrCreateReference)
				assert.ErrorIs(t, err, taxonomy.ErrReferenceNotCreated)
			},
		},
		{
			name: `happy flow`,
			args: args{
				termID:     6,
				namespace:  `emma`,
				entitiesID: []model.EntityID{`smart`, `smart`, `beautiful`, `pretty`, `clever`, `strong`},
			},
			config: func() *reference.Config {
				ns := mock.Mock[taxonomy.Namespace]()
				mock.When(ns.GetByName(mock.Any[context.Context](), mock.Exact[string](`emma`))).
					ThenReturn(&model.Namespace{ID: 66}, nil)

				trm := mock.Mock[taxonomy.Term]()
				mock.When(trm.GetByID(mock.Any[context.Context](), mock.Exact[uint64](6))).
					ThenReturn(&model.Term{}, nil)

				ref := mock.Mock[repository.Reference]()
				mock.When(ref.Set(mock.Any[context.Context](), mock.Any[[]*repository.ReferenceModel]()...)).
					ThenReturn(nil)

				return &reference.Config{
					NamespaceService:    ns,
					TermService:         trm,
					ReferenceRepository: ref,
				}
			},
			check: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.SetUp(t)
			var ctx = context.Background()
			var config = tt.config()
			config.Logger = zap.NewNop()
			r := reference.New(config)
			tt.check(t, r.Create(ctx, tt.args.termID, tt.args.namespace, tt.args.entitiesID...))
		})
	}
}

func TestService_Delete(t *testing.T) {
	type args struct {
		ctx        context.Context
		termID     uint64
		namespace  string
		entitiesID []model.EntityID
	}
	tests := []struct {
		name   string
		args   args
		check  func(t *testing.T, err error)
		config func() *reference.Config
	}{
		{
			name: `namespace not found`,
			args: args{
				termID:    6,
				namespace: `kleo`,
			},
			config: func() *reference.Config {
				ns := mock.Mock[taxonomy.Namespace]()
				mock.When(ns.GetByName(mock.Any[context.Context](), mock.Exact[string](`kleo`))).
					ThenReturn(nil, io.EOF)

				return &reference.Config{
					NamespaceService: ns,
				}
			},
			check: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, taxonomy.ErrNamespaceNotFound)
				assert.ErrorIs(t, err, io.EOF)
				assert.ErrorContains(t, err, `kleo`)
			},
		},
		{
			name: `delete error`,
			args: args{
				termID:     6,
				namespace:  `kleo`,
				entitiesID: []model.EntityID{`qwe`, `rty`},
			},
			config: func() *reference.Config {
				ns := mock.Mock[taxonomy.Namespace]()
				mock.When(ns.GetByName(mock.Any[context.Context](), mock.Exact[string](`kleo`))).
					ThenReturn(&model.Namespace{ID: 2}, nil)

				ref := mock.Mock[repository.Reference]()
				mock.When(ref.Delete(mock.Any[context.Context](), mock.Any[*repository.ReferenceFilter]())).
					ThenReturn(io.EOF)

				return &reference.Config{
					NamespaceService:    ns,
					ReferenceRepository: ref,
				}
			},
			check: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, taxonomy.ErrReferenceNotRemoved)
				assert.ErrorIs(t, err, io.EOF)
			},
		},
		{
			name: `no error`,
			args: args{
				termID:     6,
				namespace:  `kleo`,
				entitiesID: []model.EntityID{`qwe`, `rty`},
			},
			config: func() *reference.Config {
				ns := mock.Mock[taxonomy.Namespace]()
				mock.When(ns.GetByName(mock.Any[context.Context](), mock.Exact[string](`kleo`))).
					ThenReturn(&model.Namespace{ID: 2}, nil)

				ref := mock.Mock[repository.Reference]()
				mock.When(ref.Delete(mock.Any[context.Context](), mock.Any[*repository.ReferenceFilter]())).
					ThenReturn(nil)

				return &reference.Config{
					NamespaceService:    ns,
					ReferenceRepository: ref,
				}
			},
			check: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.SetUp(t)
			var ctx = context.Background()
			var config = tt.config()
			config.Logger = zap.NewNop()
			r := reference.New(config)
			tt.check(t, r.Delete(ctx, tt.args.termID, tt.args.namespace, tt.args.entitiesID...))
		})
	}
}

func TestService_Get(t *testing.T) {
	tests := []struct {
		name   string
		filter *model.ReferenceFilter
		check  func(t *testing.T, entities []*model.Reference, err error)
		config func() *reference.Config
	}{
		{
			name: `namespace not found`,
			filter: &model.ReferenceFilter{
				Namespace: []string{`ara`, `pidor`, `eba`},
			},
			config: func() *reference.Config {
				ns := mock.Mock[taxonomy.Namespace]()
				mock.When(ns.GetByName(mock.Any[context.Context](), mock.Exact[string](`ara`))).
					ThenReturn(&model.Namespace{ID: 2}, nil)
				mock.When(ns.GetByName(mock.Any[context.Context](), mock.Exact[string](`pidor`))).
					ThenReturn(nil, repository.ErrFindNamespace)

				return &reference.Config{
					NamespaceService: ns,
				}
			},
			check: func(t *testing.T, entities []*model.Reference, err error) {
				assert.Empty(t, entities)
				assert.ErrorIs(t, err, repository.ErrFindNamespace)
				assert.ErrorIs(t, err, taxonomy.ErrNamespaceNotFound)
			},
		},
		{
			name: `get references error`,
			filter: &model.ReferenceFilter{
				Namespace: []string{`ara`, `pidor`},
			},
			config: func() *reference.Config {
				ns := mock.Mock[taxonomy.Namespace]()
				mock.When(ns.GetByName(mock.Any[context.Context](), mock.Exact[string](`ara`))).
					ThenReturn(&model.Namespace{ID: 2}, nil)
				mock.When(ns.GetByName(mock.Any[context.Context](), mock.Exact[string](`pidor`))).
					ThenReturn(&model.Namespace{ID: 3}, nil)

				ref := mock.Mock[repository.Reference]()
				mock.When(ref.Get(mock.Any[context.Context](), mock.Any[*repository.ReferenceFilter]())).
					ThenReturn(nil, io.EOF)

				return &reference.Config{
					NamespaceService:    ns,
					ReferenceRepository: ref,
				}
			},
			check: func(t *testing.T, entities []*model.Reference, err error) {
				assert.Empty(t, entities)
				assert.ErrorIs(t, err, io.EOF)
			},
		},
		{
			name: `happy flow`,
			filter: &model.ReferenceFilter{
				Namespace: []string{`ara`, `pidor`},
				TermID:    [][]uint64{{99, 66}, {88}},
			},
			config: func() *reference.Config {
				ns := mock.Mock[taxonomy.Namespace]()
				mock.When(ns.GetByName(mock.Any[context.Context](), mock.Exact[string](`ara`))).
					ThenReturn(&model.Namespace{ID: 2}, nil)
				mock.When(ns.GetByName(mock.Any[context.Context](), mock.Exact[string](`pidor`))).
					ThenReturn(&model.Namespace{ID: 3}, nil)

				ref := mock.Mock[repository.Reference]()
				mock.When(ref.Get(mock.Any[context.Context](), mock.Any[*repository.ReferenceFilter]())).
					ThenAnswer(func(args []any) []any {
						var filter = args[1].(*repository.ReferenceFilter)
						var fkr = faker.New()
						var output = make([]*repository.ReferenceModel, 0)

						for _, ns := range filter.NamespaceID {
							for _, termID := range filter.TermID {
								for _, trmID := range termID {
									output = append(output, &repository.ReferenceModel{
										ID:          fkr.UInt64(),
										NamespaceID: ns,
										TermID:      trmID,
										EntityID:    model.EntityID(fkr.Beer().Name()),
									})
								}
							}
						}

						return []any{output, nil}
					})

				return &reference.Config{
					NamespaceService:    ns,
					ReferenceRepository: ref,
				}
			},
			check: func(t *testing.T, entities []*model.Reference, err error) {
				assert.NoError(t, err)
				assert.Len(t, entities, 6)
				assert.Equal(t, 66, int(entities[1].TermID))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.SetUp(t)
			var ctx = context.Background()
			var config = tt.config()
			config.Logger = zap.NewNop()
			r := reference.New(config)
			entities, err := r.Get(ctx, tt.filter)
			tt.check(t, entities, err)
		})
	}
}

//
//func TestService_GetTermsByEntities(t1 *testing.T) {
//	type fields struct {
//		log                 *zap.Logger
//		namespaceService    taxonomy.Namespace
//		referenceRepository repository.Reference
//		termService         taxonomy.Term
//	}
//	type args struct {
//		ctx       context.Context
//		namespace string
//		entities  []model.EntityID
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    []model.Term
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t1.Run(tt.name, func(t1 *testing.T) {
//			t := &Service{
//				log:                 tt.fields.log,
//				namespaceService:    tt.fields.namespaceService,
//				referenceRepository: tt.fields.referenceRepository,
//				termService:         tt.fields.termService,
//			}
//			got, err := t.GetTermsByEntities(tt.args.ctx, tt.args.namespace, tt.args.entities...)
//			if (err != nil) != tt.wantErr {
//				t1.Errorf("GetTermsByEntities() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t1.Errorf("GetTermsByEntities() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestTermService_GetReferences(t *testing.T) {
//	//	t.Parallel()
//	//
//	//	tests := []struct {
//	//		name                      string
//	//		termGroups                [][]uint
//	//		ReferenceGetReturns       func() ([]model.Reference, error)
//	//		NamespaceGetByNameReturns func() (model.NamespaceID, error)
//	//		want                      func(t assert.TestingT, references []model.Reference)
//	//		wantErr                   assert.ErrorAssertionFunc
//	//	}{
//	//		{
//	//			name:       `no references`,
//	//			termGroups: [][]uint{{34}, {92, 96}},
//	//			NamespaceGetByNameReturns: func() (model.NamespaceID, error) {
//	//				return model.NamespaceID{ID: 11}, nil
//	//			},
//	//			ReferenceGetReturns: func() ([]model.Reference, error) {
//	//				return nil, nil
//	//			},
//	//			want: func(t assert.TestingT, references []model.Reference) {
//	//				assert.Len(t, references, 0)
//	//			},
//	//			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
//	//				return assert.NoError(t, err)
//	//			},
//	//		},
//	//		{
//	//			name:       `namespace not found`,
//	//			termGroups: [][]uint{{34}, {92, 96}},
//	//			NamespaceGetByNameReturns: func() (model.NamespaceID, error) {
//	//				return model.NamespaceID{ID: 11}, errors.New(``)
//	//			},
//	//			ReferenceGetReturns: func() ([]model.Reference, error) {
//	//				return nil, nil
//	//			},
//	//			want: func(t assert.TestingT, references []model.Reference) {
//	//				assert.Len(t, references, 0)
//	//			},
//	//			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
//	//				return assert.ErrorIs(t, err, taxonomy.ErrTermNamespaceNotFound)
//	//			},
//	//		},
//	//	}
//	//
//	//	// Create categories "RAM", "CPU", "Display size"
//	//	// Should receive all laptops that has: "RAM" (512 or 1024) and "CPU" (2.8 or 3.2) and "Display size"
//	//	// (between 13 and 15)
//	//
//	//	for _, tt := range tests {
//	//		t.Run(tt.name, func(t *testing.T) {
//	//			t.Parallel()
//	//			referencerepo := mockrepository.NewReference(t)
//	//			referencerepo.On(`Get`, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
//	//				Return(tt.ReferenceGetReturns()).Maybe()
//	//
//	//			namespaceservice := mockservice.NewNamespace(t)
//	//			namespaceservice.On(`GetByName`, mock.Anything, mock.Anything).Return(tt.NamespaceGetByNameReturns())
//	//
//	//			termService := reference.New(&reference.Config{
//	//				NamespaceService:    namespaceservice,
//	//				ReferenceRepository: referencerepo,
//	//				Logger:              zap.NewNop(),
//	//			})
//	//
//	//			references, err := termService.GetReferences(context.TODO(), &model.EntityFilter{
//	//				TermID:      tt.termGroups,
//	//				NamespaceID: []string{`any`},
//	//			})
//	//			tt.wantErr(t, err)
//	//			tt.want(t, references)
//	//		})
//	//	}
//}
