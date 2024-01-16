package repository_test

import (
	"context"
	"github.com/dmalykh/taxonomy/internal/repository/entgo/ent"
	"github.com/dmalykh/taxonomy/internal/repository/entgo/ent/enttest"
	repo "github.com/dmalykh/taxonomy/internal/repository/entgo/repository"
	"github.com/dmalykh/taxonomy/taxonomy/model"
	"github.com/dmalykh/taxonomy/taxonomy/repository"
	"github.com/samber/lo"
	"testing"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ReferenceTestSuite struct {
	suite.Suite
	client *ent.Client
	faker  faker.Faker
}

func (suite *ReferenceTestSuite) SetupTest() {
	suite.client = enttest.Open(suite.T(), "sqlite3", ":memory:?_fk=1", []enttest.Option{
		enttest.WithOptions(ent.Log(suite.T().Log)),
	}...).Debug()
	suite.faker = faker.New()
}

func (suite *ReferenceTestSuite) TearDownTest() {
	assert.NoError(suite.T(), suite.client.Close())
}

func (suite *ReferenceTestSuite) mockTerm(ctx context.Context, vocabularyID uint64) *ent.Term {
	return suite.client.Term.Create().
		SetName(suite.faker.RandomStringWithLength(suite.faker.IntBetween(1, 9999))).
		SetTitle(suite.faker.Beer().Name()).
		AddVocabularyIDs(vocabularyID).
		SaveX(ctx)
}

func (suite *ReferenceTestSuite) mockVocabulary(ctx context.Context, parentID *uint64) *ent.Vocabulary {
	return suite.client.Vocabulary.Create().
		SetName(suite.faker.RandomStringWithLength(suite.faker.IntBetween(1, 9999))).
		SetTitle(suite.faker.Company().Name()).
		SetNillableParentID(parentID).
		SaveX(ctx)
}

func (suite *ReferenceTestSuite) mockNamespace(ctx context.Context) *ent.Namespace {
	return suite.client.Namespace.Create().
		SetName(suite.faker.RandomStringWithLength(suite.faker.IntBetween(1, 9999))).
		SaveX(ctx)
}

func (suite *ReferenceTestSuite) mockReference(ctx context.Context, termID, namespaceID uint64, entityID string) *ent.Reference {
	return suite.client.Reference.Create().
		SetTermID(termID).
		SetNamespaceID(namespaceID).
		SetEntityID(entityID).
		SaveX(ctx)
}

// Generate mock data for deletion
func (suite *ReferenceTestSuite) generate(count int) ([]uint64, []uint64, []model.EntityID) {
	var (
		ctx                         = context.Background()
		terms, namespaces, entities = make([]uint64, count), make([]uint64, count), make([]model.EntityID, count)
	)

	for i := 0; i < count; i++ {
		var (
			vocabulary = suite.mockVocabulary(ctx, nil)
			term       = suite.mockTerm(ctx, vocabulary.ID)
			namespace  = suite.mockNamespace(ctx)
			entityID   = suite.faker.UUID().V4()
		)

		suite.mockReference(ctx, term.ID, namespace.ID, entityID)
		terms[i], namespaces[i], entities[i] = term.ID, namespace.ID, model.EntityID(entityID)
	}

	return terms, namespaces, entities
}

func (suite *ReferenceTestSuite) TestCreate() {
	tests := []struct {
		name       string
		references func(ctx context.Context, client *ent.Client) []*repository.ReferenceModel
		err        error
		want       int
	}{
		{
			name: `error on empty values`,
			references: func(ctx context.Context, client *ent.Client) []*repository.ReferenceModel {
				return []*repository.ReferenceModel{
					{},
				}
			},
			err:  repository.ErrCreateReference,
			want: 0,
		},
		{
			name: `ok`,
			references: func(ctx context.Context, client *ent.Client) []*repository.ReferenceModel {
				vocabulary := suite.mockVocabulary(ctx, nil)

				return []*repository.ReferenceModel{
					{
						TermID:      suite.mockTerm(ctx, vocabulary.ID).ID,
						NamespaceID: suite.mockNamespace(ctx).ID,
						EntityID:    `140`,
					},
					{
						TermID:      suite.mockTerm(ctx, vocabulary.ID).ID,
						NamespaceID: suite.mockNamespace(ctx).ID,
						EntityID:    `99999999`,
					},
				}
			},
			err:  nil,
			want: 2,
		},
		{
			name: `error, one of references broken`,
			references: func(ctx context.Context, client *ent.Client) []*repository.ReferenceModel {
				vocabulary := client.Vocabulary.Create().SetName(suite.faker.Company().Suffix()).
					SetTitle(suite.faker.Company().Name()).SaveX(ctx)

				return []*repository.ReferenceModel{
					{
						TermID:      suite.mockTerm(ctx, vocabulary.ID).ID,
						NamespaceID: suite.mockNamespace(ctx).ID,
						EntityID:    `140`,
					},
					{
						NamespaceID: suite.mockNamespace(ctx).ID,
						EntityID:    `99999999`,
					},
				}
			},
			err:  repository.ErrCreateReference,
			want: 0,
		},
		{
			name: `no error same terms and namespaces`,
			references: func(ctx context.Context, client *ent.Client) []*repository.ReferenceModel {
				vocabulary := client.Vocabulary.Create().SetName(suite.faker.Company().Suffix()).
					SetTitle(suite.faker.Company().Name()).SaveX(ctx)
				term := suite.mockTerm(ctx, vocabulary.ID)
				namespace := suite.mockNamespace(ctx)

				return []*repository.ReferenceModel{
					{
						TermID:      term.ID,
						NamespaceID: namespace.ID,
						EntityID:    `140`,
					},
					{
						TermID:      term.ID,
						NamespaceID: namespace.ID,
						EntityID:    `99999999`,
					},
				}
			},
			err:  nil,
			want: 2,
		},
		{
			name: `duplicate records error`,
			references: func(ctx context.Context, client *ent.Client) []*repository.ReferenceModel {
				vocabulary := client.Vocabulary.Create().SetName(suite.faker.Company().Suffix()).
					SetTitle(suite.faker.Company().Name()).SaveX(ctx)
				term := suite.mockTerm(ctx, vocabulary.ID)
				namespace := suite.mockNamespace(ctx)

				return []*repository.ReferenceModel{
					{
						TermID:      term.ID,
						NamespaceID: namespace.ID,
						EntityID:    `222`,
					},
					{
						TermID:      term.ID,
						NamespaceID: namespace.ID,
						EntityID:    `222`,
					},
				}
			},
			err:  repository.ErrCreateReference,
			want: 0,
		},
	}

	for _, tt := range tests {
		suite.TearDownTest()
		suite.SetupTest()
		suite.Run(tt.name, func() {
			ctx := context.TODO()
			{
				rel := repo.NewReference(suite.client.Reference)
				err := rel.Create(ctx, tt.references(ctx, suite.client)...)
				assert.ErrorIs(suite.T(), err, tt.err)
			}
			{
				count, err := suite.client.Reference.Query().Count(ctx)
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tt.want, count)
			}
		})
	}
}

func (suite *ReferenceTestSuite) TestDelete() {
	tests := []struct {
		name   string
		filter func() *repository.ReferenceFilter
		check  func(t assert.TestingT, err error)
	}{
		{
			name: `without namespace error`,
			filter: func() *repository.ReferenceFilter {
				return &repository.ReferenceFilter{
					EntityID: []model.EntityID{`22`, `33`},
				}
			},
			check: func(t assert.TestingT, err error) {
				assert.ErrorIs(t, err, repository.ErrWithoutNamespace)
			},
		},
		{
			name: `remove half of terms only`,
			filter: func() *repository.ReferenceFilter {
				terms, ns, _ := suite.generate(100)
				return &repository.ReferenceFilter{
					NamespaceID: ns,
					TermID:      [][]uint64{terms[:50]},
				}
			},
			check: func(t assert.TestingT, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 50, suite.client.Reference.Query().CountX(context.TODO()))
			},
		},
		{
			name: `remove half of namespaces only`,
			filter: func() *repository.ReferenceFilter {
				_, ns, _ := suite.generate(100)
				return &repository.ReferenceFilter{
					NamespaceID: ns[:50],
				}
			},
			check: func(t assert.TestingT, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 50, suite.client.Reference.Query().CountX(context.TODO()))
			},
		},
		{
			name: `remove half of entities only`,
			filter: func() *repository.ReferenceFilter {
				_, ns, entities := suite.generate(100)
				return &repository.ReferenceFilter{
					NamespaceID: ns,
					EntityID:    entities[:50],
				}
			},
			check: func(t assert.TestingT, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 50, suite.client.Reference.Query().CountX(context.TODO()))
			},
		},
	}

	for _, tt := range tests {
		suite.TearDownTest()
		suite.SetupTest()
		suite.Run(tt.name, func() {
			ctx := context.TODO()
			rel := repo.NewReference(suite.client.Reference)
			filter := tt.filter()
			tt.check(suite.T(), rel.Delete(ctx, filter))
		})
	}
}

func (suite *ReferenceTestSuite) TestGet() {
	tests := []struct {
		name  string
		get   func() ([][]uint64, []uint64, []model.EntityID)
		check func(references []repository.ReferenceModel)
		err   func(error, ...interface{}) bool
	}{
		{
			`get all entities`,
			func() ([][]uint64, []uint64, []model.EntityID) {
				_, ns, _ := suite.generate(100)

				return nil, ns, nil
			},
			func(references []repository.ReferenceModel) {
				suite.Len(references, 100)
			},
			func(err error, i ...interface{}) bool {
				return suite.NoError(err)
			},
		},
		{
			`entity without namespace error`,
			func() ([][]uint64, []uint64, []model.EntityID) {
				return nil, nil, []model.EntityID{`22`, `33`}
			},
			func(references []repository.ReferenceModel) {
				suite.Empty(references)
			},
			func(err error, i ...interface{}) bool {
				return suite.ErrorIs(err, repository.ErrWithoutNamespace)
			},
		},
		{
			`get terms by complex condition`,
			func() ([][]uint64, []uint64, []model.EntityID) {
				//terms, ns, _ := suite.generate(1000)
				var ctx = context.Background()
				ns := suite.mockNamespace(ctx)

				vocRAM := suite.mockVocabulary(ctx, nil)
				t32 := suite.mockTerm(ctx, vocRAM.ID)
				t256 := suite.mockTerm(ctx, vocRAM.ID)

				vocMatrix := suite.mockVocabulary(ctx, nil)
				tOLED := suite.mockTerm(ctx, vocMatrix.ID)
				tIPS := suite.mockTerm(ctx, vocMatrix.ID)
				tTFT := suite.mockTerm(ctx, vocMatrix.ID)

				suite.mockReference(ctx, t32.ID, ns.ID, `dellX`)
				suite.mockReference(ctx, tOLED.ID, ns.ID, `dellX`)

				suite.mockReference(ctx, t32.ID, ns.ID, `dellY`)
				suite.mockReference(ctx, tIPS.ID, ns.ID, `dellY`)

				suite.mockReference(ctx, t256.ID, ns.ID, `dellZ`)
				suite.mockReference(ctx, tTFT.ID, ns.ID, `dellZ`)

				suite.mockReference(ctx, t256.ID, ns.ID, `dellM`)
				suite.mockReference(ctx, tIPS.ID, ns.ID, `dellM`)

				suite.mockReference(ctx, t256.ID, ns.ID, `asusX`)
				suite.mockReference(ctx, tOLED.ID, ns.ID, `asusX`)

				suite.mockReference(ctx, t32.ID, ns.ID, `asusY`)
				suite.mockReference(ctx, tIPS.ID, ns.ID, `asusY`)

				// ALl laptops with 256 RAM and OLED or IPS matrix
				return [][]uint64{{t256.ID}, {tOLED.ID, tIPS.ID}}, []uint64{ns.ID}, nil
			},
			func(references []repository.ReferenceModel) {
				// Should be dellM and asusX (4 terms
				suite.Len(references, 4)
				suite.Len(lo.Filter[repository.ReferenceModel](references, func(item repository.ReferenceModel, index int) bool {
					return item.EntityID != `dellM` && item.EntityID != `asusX`
				}), 0)
			},
			func(err error, i ...interface{}) bool {
				return suite.NoError(err)
			},
		},
		{
			`get half of namespaces only`,
			func() ([][]uint64, []uint64, []model.EntityID) {
				_, namespaces, _ := suite.generate(100)

				return nil, namespaces[:50], nil
			},
			func(references []repository.ReferenceModel) {
				suite.Len(references, 50)
			},
			func(err error, i ...interface{}) bool {
				return suite.NoError(err)
			},
		},
		{
			`get half of entities only`,
			func() ([][]uint64, []uint64, []model.EntityID) {
				_, namespaces, entities := suite.generate(100)

				return nil, namespaces, entities[:50]
			},
			func(references []repository.ReferenceModel) {
				suite.Len(references, 50)
			},
			func(err error, i ...interface{}) bool {
				return suite.NoError(err)
			},
		},
	}

	for _, tt := range tests {
		suite.TearDownTest()
		suite.SetupTest()
		suite.Run(tt.name, func() {
			ctx := context.TODO()
			rel := repo.NewReference(suite.client.Reference)
			termIds, namespaceIds, entityIds := tt.get()
			references, err := rel.Get(ctx, &repository.ReferenceFilter{
				TermID:      termIds,
				NamespaceID: namespaceIds,
				EntityID:    entityIds,
			})
			require.True(suite.T(), tt.err(err))
			tt.check(references)
		})
	}
}

func TestReferenceTestSuite(t *testing.T) {
	suite.Run(t, new(ReferenceTestSuite))
}
