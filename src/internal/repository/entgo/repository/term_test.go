package repository_test

import (
	"context"
	"github.com/dmalykh/taxonomy/internal/repository/entgo/ent"
	"github.com/dmalykh/taxonomy/internal/repository/entgo/ent/enttest"
	repo "github.com/dmalykh/taxonomy/internal/repository/entgo/repository"
	"github.com/dmalykh/taxonomy/taxonomy/model"
	"github.com/dmalykh/taxonomy/taxonomy/repository"
	"testing"

	"github.com/jaswdr/faker"
	suitetest "github.com/stretchr/testify/suite"
)

type TestTermOperations struct {
	suitetest.Suite
	client *ent.Client
}

func (suite *TestTermOperations) SetupTest() {
	suite.client = enttest.Open(suite.T(), "sqlite3", ":memory:?_fk=1", []enttest.Option{
		enttest.WithOptions(ent.Log(suite.T().Log)),
	}...) // .Debug()
}

func (suite *TestTermOperations) TearDownTest() {

	suite.client.Close()
}

func (suite *TestTermOperations) TestTerm_Create() {

	faker := faker.New()

	tests := []struct {
		name     string
		prepare  func()
		data     model.TermData
		checkErr func(error, ...interface{}) // returns continue
	}{
		{
			`ok`,
			func() {
				suite.client.Vocabulary.Create().SetName(faker.Beer().Name()).SetTitle(``).SaveX(context.TODO())
			},
			model.TermData{
				Name:         faker.Beer().Name(),
				Title:        faker.Beer().Name(),
				VocabularyID: []uint64{1},
			},
			func(err error, i ...interface{}) {
				suite.NoError(err)
			},
		},
		{
			`title empty, no error`,
			func() {
				suite.client.Vocabulary.Create().SetName(faker.Beer().Name()).SetTitle(``).SaveX(context.TODO())
			},
			model.TermData{
				Name:         faker.Beer().Name(),
				VocabularyID: []uint64{1},
			},
			func(err error, i ...interface{}) {
				suite.NoError(err)
			},
		},
		{
			`name empty, error`,
			func() {
				suite.client.Vocabulary.Create().SetName(faker.Beer().Name()).SetTitle(``).SaveX(context.TODO())
			},
			model.TermData{
				Title:        faker.Beer().Name(),
				VocabularyID: []uint64{1},
			},
			func(err error, i ...interface{}) {
				suite.ErrorIs(err, repository.ErrCreateTerm)
			},
		},
		{
			`vocabulary empty, error`,
			func() {},
			model.TermData{
				Name:  faker.Beer().Name(),
				Title: faker.Beer().Name(),
			},
			func(err error, i ...interface{}) {
				suite.ErrorIs(err, repository.ErrCreateTerm)
			},
		},
		{
			`duplicated names, unique categories`,
			func() {
				suite.client.Vocabulary.Create().SetName(faker.Beer().Name()).SetTitle(``).SaveX(context.TODO())
				suite.client.Vocabulary.Create().SetName(faker.Company().Name()).SetTitle(``).SaveX(context.TODO())
				suite.client.Term.Create().SetName(`sowa`).AddVocabularyIDs(1).SaveX(context.TODO())
			},
			model.TermData{
				Name:         `sowa`,
				VocabularyID: []uint64{2},
			},
			func(err error, i ...interface{}) {
				suite.NoError(err)
			},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.TearDownTest()
			suite.SetupTest()
			ctx := context.TODO()
			termClient := repo.NewTerm(suite.client.Term)

			tt.prepare()
			returned, err := termClient.Create(ctx, &tt.data)
			tt.checkErr(err)
			if err != nil {
				return
			}
			{
				suite.Equal(tt.data.Name, returned.Data.Name)
				suite.Equal(tt.data.Title, returned.Data.Title)
				suite.Equal(tt.data.Description, returned.Data.Description)
				suite.EqualValues(tt.data.VocabularyID, returned.Data.VocabularyID)
			}

			{
				got, err := termClient.Get(ctx, &repository.TermFilter{ID: []uint64{returned.ID}})
				suite.NoError(err)
				suite.Len(got, 1)
				suite.Equal(tt.data.Name, got[0].Data.Name)
				suite.Equal(tt.data.Title, got[0].Data.Title)
				suite.Equal(tt.data.Description, got[0].Data.Description)
				suite.Equal(tt.data.VocabularyID, got[0].Data.VocabularyID)
			}
		})
	}
}

func TestTermOperationsSuite(t *testing.T) {
	suitetest.Run(t, new(TestTermOperations))
}
