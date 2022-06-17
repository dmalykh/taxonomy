package repository_test

import (
	"context"
	"testing"

	"github.com/dmalykh/tagservice/repository/entgo/ent"
	"github.com/dmalykh/tagservice/repository/entgo/ent/enttest"
	repo "github.com/dmalykh/tagservice/repository/entgo/repository"
	"github.com/dmalykh/tagservice/tagservice/model"
	"github.com/dmalykh/tagservice/tagservice/repository"
	"github.com/jaswdr/faker"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RelationTestSuite struct {
	suite.Suite
	client *ent.Client
	faker  faker.Faker
}

func (suite *RelationTestSuite) SetupTest() {
	suite.client = enttest.Open(suite.T(), "sqlite3", ":memory:?_fk=1", []enttest.Option{
		enttest.WithOptions(ent.Log(suite.T().Log)),
	}...) // .Debug()
	suite.faker = faker.New()
}

func (suite *RelationTestSuite) TearDownTest() {
	assert.NoError(suite.T(), suite.client.Close())
}

func (suite *RelationTestSuite) mockTag(ctx context.Context, categoryID int) *ent.Tag {
	return suite.client.Tag.Create().
		SetName(suite.faker.RandomStringWithLength(suite.faker.IntBetween(1, 9999))).
		SetTitle(suite.faker.Beer().Name()).
		SetCategoryID(categoryID).
		SaveX(ctx)
}

func (suite *RelationTestSuite) mockCategory(ctx context.Context, parentID *int) *ent.Category {
	return suite.client.Category.Create().
		SetName(suite.faker.RandomStringWithLength(suite.faker.IntBetween(1, 9999))).
		SetTitle(suite.faker.Company().Name()).
		SetNillableParentID(parentID).
		SaveX(ctx)
}

func (suite *RelationTestSuite) mockNamespace(ctx context.Context) *ent.Namespace {
	return suite.client.Namespace.Create().
		SetName(suite.faker.RandomStringWithLength(suite.faker.IntBetween(1, 9999))).
		SaveX(ctx)
}

func (suite *RelationTestSuite) mockRelation(ctx context.Context, tagID, namespaceID, entityID int) *ent.Relation {
	return suite.client.Relation.Create().
		SetTagID(tagID).
		SetNamespaceID(namespaceID).
		SetEntityID(entityID).
		SaveX(ctx)
}

// Generate mock data for deletion
//goland:noinspection GoContextTodo
func (suite *RelationTestSuite) generate(count int) ([]uint, []uint, []uint) {
	var (
		ctx                        = context.TODO()
		tags, namespaces, entities = make([]uint, count), make([]uint, count), make([]uint, count)
	)

	for i := 0; i < count; i++ {
		var (
			category  = suite.mockCategory(ctx, nil)
			tag       = suite.mockTag(ctx, category.ID)
			namespace = suite.mockNamespace(ctx)
			entityID  = suite.faker.Int()
		)

		suite.mockRelation(ctx, tag.ID, namespace.ID, entityID)
		tags[i], namespaces[i], entities[i] = uint(tag.ID), uint(namespace.ID), uint(entityID)
	}

	return tags, namespaces, entities
}

//goland:noinspection GoContextTodo
func (suite *RelationTestSuite) TestCreate() {
	tests := []struct {
		name      string
		relations func(ctx context.Context, client *ent.Client) []*model.Relation
		err       error
		want      int
	}{
		{
			name: `error on empty values`,
			relations: func(ctx context.Context, client *ent.Client) []*model.Relation {
				return []*model.Relation{
					{},
				}
			},
			err:  repository.ErrCreateRelation,
			want: 0,
		},
		{
			name: `ok`,
			relations: func(ctx context.Context, client *ent.Client) []*model.Relation {
				category := suite.mockCategory(ctx, nil)

				return []*model.Relation{
					{
						TagID:       uint(suite.mockTag(ctx, category.ID).ID),
						NamespaceID: uint(suite.mockNamespace(ctx).ID),
						EntityID:    140,
					},
					{
						TagID:       uint(suite.mockTag(ctx, category.ID).ID),
						NamespaceID: uint(suite.mockNamespace(ctx).ID),
						EntityID:    99999999,
					},
				}
			},
			err:  nil,
			want: 2,
		},
		{
			name: `error, one of relations broken`,
			relations: func(ctx context.Context, client *ent.Client) []*model.Relation {
				category := client.Category.Create().SetName(suite.faker.Company().Suffix()).
					SetTitle(suite.faker.Company().Name()).SaveX(ctx)

				return []*model.Relation{
					{
						TagID:       uint(suite.mockTag(ctx, category.ID).ID),
						NamespaceID: uint(suite.mockNamespace(ctx).ID),
						EntityID:    140,
					},
					{
						NamespaceID: uint(suite.mockNamespace(ctx).ID),
						EntityID:    99999999,
					},
				}
			},
			err:  repository.ErrCreateRelation,
			want: 0,
		},
		{
			name: `no error same tags and namespaces`,
			relations: func(ctx context.Context, client *ent.Client) []*model.Relation {
				category := client.Category.Create().SetName(suite.faker.Company().Suffix()).
					SetTitle(suite.faker.Company().Name()).SaveX(ctx)
				tag := suite.mockTag(ctx, category.ID)
				namespace := suite.mockNamespace(ctx)

				return []*model.Relation{
					{
						TagID:       uint(tag.ID),
						NamespaceID: uint(namespace.ID),
						EntityID:    140,
					},
					{
						TagID:       uint(tag.ID),
						NamespaceID: uint(namespace.ID),
						EntityID:    99999999,
					},
				}
			},
			err:  nil,
			want: 2,
		},
		{
			name: `duplicate records error`,
			relations: func(ctx context.Context, client *ent.Client) []*model.Relation {
				category := client.Category.Create().SetName(suite.faker.Company().Suffix()).
					SetTitle(suite.faker.Company().Name()).SaveX(ctx)
				tag := suite.mockTag(ctx, category.ID)
				namespace := suite.mockNamespace(ctx)

				return []*model.Relation{
					{
						TagID:       uint(tag.ID),
						NamespaceID: uint(namespace.ID),
						EntityID:    222,
					},
					{
						TagID:       uint(tag.ID),
						NamespaceID: uint(namespace.ID),
						EntityID:    222,
					},
				}
			},
			err:  nil, // @TODO make a bug in ent. CreateBulk doesn't return error, records just doesn't added.
			want: 0,
		},
	}

	for _, tt := range tests {
		suite.TearDownTest()
		suite.SetupTest()
		suite.Run(tt.name, func() {
			ctx := context.TODO()
			{
				rel := repo.NewRelation(suite.client.Relation)
				err := rel.Create(ctx, tt.relations(ctx, suite.client)...)
				assert.ErrorIs(suite.T(), err, tt.err)
			}
			{
				count, err := suite.client.Relation.Query().Count(ctx)
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tt.want, count)
			}
		})
	}
}

//goland:noinspection GoContextTodo,GoContextTodo,GoContextTodo,GoContextTodo
func (suite *RelationTestSuite) TestDelete() {
	tests := []struct {
		name   string
		delete func() ([]uint, []uint, []uint)
		check  func(t assert.TestingT)
		err    assert.ErrorAssertionFunc
	}{
		{
			`entity without namespace error`,
			func() ([]uint, []uint, []uint) {
				return nil, nil, []uint{22, 33}
			},
			func(t assert.TestingT) {
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, repository.ErrEntityWithoutNamespace)
			},
		},
		{
			`remove half of tags only`,
			func() ([]uint, []uint, []uint) {
				tags, _, _ := suite.generate(100)

				return tags[:50], nil, nil
			},
			func(t assert.TestingT) {
				assert.Equal(t, 50, suite.client.Relation.Query().CountX(context.TODO()))
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			`remove half of namespaces only`,
			func() ([]uint, []uint, []uint) {
				_, namespaces, _ := suite.generate(100)

				return nil, namespaces[:50], nil
			},
			func(t assert.TestingT) {
				assert.Equal(t, 50, suite.client.Relation.Query().CountX(context.TODO()))
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			`remove half of entities only`,
			func() ([]uint, []uint, []uint) {
				_, namespaces, entities := suite.generate(100)

				return nil, namespaces, entities[:50]
			},
			func(t assert.TestingT) {
				assert.Equal(t, 50, suite.client.Relation.Query().CountX(context.TODO()))
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		suite.TearDownTest()
		suite.SetupTest()
		suite.Run(tt.name, func() {
			ctx := context.TODO()
			rel := repo.NewRelation(suite.client.Relation)
			tagIds, namespaceIds, entityIds := tt.delete()
			require.True(suite.T(), tt.err(suite.T(), rel.Delete(ctx, tagIds, namespaceIds, entityIds)))
			tt.check(suite.T())
		})
	}
}

//goland:noinspection GoContextTodo
func (suite *RelationTestSuite) TestGet() {
	tests := []struct {
		name  string
		get   func() ([][]uint, []uint, []uint)
		check func(relations []model.Relation)
		err   func(error, ...interface{}) bool
	}{
		{
			`get all entities`,
			func() ([][]uint, []uint, []uint) {
				suite.generate(100)

				return nil, nil, nil
			},
			func(relations []model.Relation) {
				suite.Len(relations, 100)
			},
			func(err error, i ...interface{}) bool {
				return suite.NoError(err)
			},
		},
		{
			`entity without namespace error`,
			func() ([][]uint, []uint, []uint) {
				return nil, nil, []uint{22, 33}
			},
			func(relations []model.Relation) {
				suite.Empty(relations)
			},
			func(err error, i ...interface{}) bool {
				return suite.ErrorIs(err, repository.ErrEntityWithoutNamespace)
			},
		},
		{
			`get half of tags only`,
			func() ([][]uint, []uint, []uint) {
				tags, _, _ := suite.generate(100)

				return [][]uint{tags[:25], tags[25:50]}, nil, nil
			},
			func(relations []model.Relation) {
				suite.Len(relations, 50)
			},
			func(err error, i ...interface{}) bool {
				return suite.NoError(err)
			},
		},
		{
			`get half of namespaces only`,
			func() ([][]uint, []uint, []uint) {
				_, namespaces, _ := suite.generate(100)

				return nil, namespaces[:50], nil
			},
			func(relations []model.Relation) {
				suite.Len(relations, 50)
			},
			func(err error, i ...interface{}) bool {
				return suite.NoError(err)
			},
		},
		{
			`get half of entities only`,
			func() ([][]uint, []uint, []uint) {
				_, namespaces, entities := suite.generate(100)

				return nil, namespaces, entities[:50]
			},
			func(relations []model.Relation) {
				suite.Len(relations, 50)
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
			rel := repo.NewRelation(suite.client.Relation)
			tagIds, namespaceIds, entityIds := tt.get()
			relations, err := rel.Get(ctx, &model.RelationFilter{
				TagID:     tagIds,
				Namespace: namespaceIds,
				EntityID:  entityIds,
			})
			require.True(suite.T(), tt.err(err))
			tt.check(relations)
		})
	}
}

func TestRelationTestSuite(t *testing.T) {
	suite.Run(t, new(RelationTestSuite))
}
