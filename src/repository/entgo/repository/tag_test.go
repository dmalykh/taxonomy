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
	suitetest "github.com/stretchr/testify/suite"
)

type TestTagOperations struct {
	suitetest.Suite
	client *ent.Client
}

func (suite *TestTagOperations) SetupTest() {
	suite.client = enttest.Open(suite.T(), "sqlite3", ":memory:?_fk=1", []enttest.Option{
		enttest.WithOptions(ent.Log(suite.T().Log)),
	}...) // .Debug()
}

func (suite *TestTagOperations) TearDownTest() {
	//goland:noinspection GoUnhandledErrorResult
	suite.client.Close()
}

//goland:noinspection GoContextTodo
func (suite *TestTagOperations) TestTag_Create() {
	//goland:noinspection GoImportUsedAsName
	faker := faker.New()

	tests := []struct {
		name     string
		prepare  func()
		data     model.TagData
		checkErr func(error, ...interface{}) // returns continue
	}{
		{
			`ok`,
			func() {
				suite.client.Category.Create().SetName(faker.Beer().Name()).SetTitle(``).SaveX(context.TODO())
			},
			model.TagData{
				Name:       faker.Beer().Name(),
				Title:      faker.Beer().Name(),
				CategoryID: 1,
			},
			func(err error, i ...interface{}) {
				suite.NoError(err)
			},
		},
		{
			`title empty, no error`,
			func() {
				suite.client.Category.Create().SetName(faker.Beer().Name()).SetTitle(``).SaveX(context.TODO())
			},
			model.TagData{
				Name:       faker.Beer().Name(),
				CategoryID: 1,
			},
			func(err error, i ...interface{}) {
				suite.NoError(err)
			},
		},
		{
			`name empty, error`,
			func() {
				suite.client.Category.Create().SetName(faker.Beer().Name()).SetTitle(``).SaveX(context.TODO())
			},
			model.TagData{
				Title:      faker.Beer().Name(),
				CategoryID: 1,
			},
			func(err error, i ...interface{}) {
				suite.ErrorIs(err, repository.ErrCreateTag)
			},
		},
		{
			`category empty, error`,
			func() {},
			model.TagData{
				Name:  faker.Beer().Name(),
				Title: faker.Beer().Name(),
			},
			func(err error, i ...interface{}) {
				suite.ErrorIs(err, repository.ErrCreateTag)
			},
		},
		{
			`duplicated names, unique categories`,
			func() {
				suite.client.Category.Create().SetName(faker.Beer().Name()).SetTitle(``).SaveX(context.TODO())
				suite.client.Category.Create().SetName(faker.Company().Name()).SetTitle(``).SaveX(context.TODO())
				suite.client.Tag.Create().SetName(`sowa`).SetCategoryID(1).SaveX(context.TODO())
			},
			model.TagData{
				Name:       `sowa`,
				CategoryID: 2,
			},
			func(err error, i ...interface{}) {
				suite.NoError(err)
			},
		},
		{
			`duplicated names and categories`,
			func() {
				suite.client.Category.Create().SetName(faker.Beer().Name()).SetTitle(``).SaveX(context.TODO())
				suite.client.Tag.Create().SetName(`sowa`).SetCategoryID(1).SaveX(context.TODO())
			},
			model.TagData{
				Name:       `sowa`,
				CategoryID: 1,
			},
			func(err error, i ...interface{}) {
				suite.ErrorIs(err, repository.ErrCreateTag)
			},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.TearDownTest()
			suite.SetupTest()
			ctx := context.TODO()
			tagClient := repo.NewTag(suite.client.Tag)

			tt.prepare()
			returned, err := tagClient.Create(ctx, &tt.data)
			tt.checkErr(err)
			if err != nil {
				return
			}
			{
				suite.Equal(tt.data.Name, returned.Data.Name)
				suite.Equal(tt.data.Title, returned.Data.Title)
				suite.Equal(tt.data.Description, returned.Data.Description)
				suite.Equal(tt.data.CategoryID, returned.Data.CategoryID)
			}

			{
				got, err := tagClient.GetByID(ctx, returned.ID)
				suite.NoError(err)
				suite.Equal(tt.data.Name, got.Data.Name)
				suite.Equal(tt.data.Title, got.Data.Title)
				suite.Equal(tt.data.Description, got.Data.Description)
				suite.Equal(tt.data.CategoryID, got.Data.CategoryID)
			}
		})
	}
}

func TestTagOperationsSuite(t *testing.T) {
	suitetest.Run(t, new(TestTagOperations))
}
