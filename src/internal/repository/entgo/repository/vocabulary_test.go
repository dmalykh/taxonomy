package repository_test

import (
	"context"
	"github.com/dmalykh/taxonomy/internal/repository/entgo/ent"
	"github.com/dmalykh/taxonomy/internal/repository/entgo/ent/enttest"
	repo "github.com/dmalykh/taxonomy/internal/repository/entgo/repository"
	"github.com/dmalykh/taxonomy/taxonomy/model"
	"github.com/dmalykh/taxonomy/taxonomy/repository"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVocabulary_Create(t *testing.T) {
	faker := faker.New()

	tests := []struct {
		name    string
		data    model.VocabularyData
		prepare func(c *ent.Client)
		wantErr assert.ErrorAssertionFunc // returns continue
		check   func(t assert.TestingT, c *ent.Client)
	}{
		{
			`ok`,
			model.VocabularyData{
				Name:  faker.Beer().Name(),
				Title: faker.Beer().Name(),
			},
			func(c *ent.Client) {},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
			func(t assert.TestingT, c *ent.Client) {},
		},
		{
			`title empty, no error`,
			model.VocabularyData{
				Name: faker.Beer().Name(),
			},
			func(c *ent.Client) {},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
			func(t assert.TestingT, c *ent.Client) {},
		},
		{
			`name empty, error`,
			model.VocabularyData{
				Title: faker.Beer().Name(),
			},
			func(c *ent.Client) {},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.ErrorIs(t, err, repository.ErrCreateVocabulary)

				return false
			},
			func(t assert.TestingT, c *ent.Client) {},
		},
		{
			`name duplicated, parentid unique, no error`,
			model.VocabularyData{
				Name:     `gogolek`,
				Title:    `kkk`,
				ParentID: pointer.ToUint64(2),
			},
			func(c *ent.Client) {
				// https://www.sqlite.org/nulls.html
				c.Vocabulary.Create().SetName(`https://www.sqlite.org/nulls.html`).SetTitle(`NULL in SQL can't be UNIQ`).SaveX(context.TODO()) // id: 1
				c.Vocabulary.Create().SetName(`gogolek`).SetTitle(`kkk`).SetParentID(1).SaveX(context.TODO())                                  // id: 2
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.NoError(t, err)

				return true
			},
			func(t assert.TestingT, c *ent.Client) {
				// Additional check. It was a bug when parent was updated too
				vocabulary := c.Vocabulary.GetX(context.TODO(), 1)
				assert.Empty(t, vocabulary.ParentID)
			},
		},
		{
			`name and parentid duplicated, error`,
			model.VocabularyData{
				Name:     `gogolek`,
				Title:    `kkk`,
				ParentID: pointer.ToUint64(1),
			},
			func(c *ent.Client) {
				// https://www.sqlite.org/nulls.html
				c.Vocabulary.Create().SetName(`https://www.sqlite.org/nulls.html`).SetTitle(`NULL in SQL can't be UNIQ`).SaveX(context.TODO())
				c.Vocabulary.Create().SetName(`gogolek`).SetTitle(`kkk`).SetParentID(1).SaveX(context.TODO())
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.NotNil(t, err)
				assert.ErrorIs(t, err, repository.ErrNotUniqueName)

				return false
			},
			func(t assert.TestingT, c *ent.Client) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			c, client := vocabularyClient(t)
			defer func() {
				tt.check(t, client)
			}()
			tt.prepare(client)

			returned, err := c.Create(ctx, &tt.data)
			if !tt.wantErr(t, err) {
				return
			}

			{
				assert.Equal(t, tt.data.Name, returned.Data.Name)
				assert.Equal(t, tt.data.Title, returned.Data.Title)
				assert.Equal(t, pointer.GetString(tt.data.Description), *returned.Data.Description)
				assert.Equal(t, tt.data.ParentID, returned.Data.ParentID)
			}

			{
				got, err := c.Get(ctx, &repository.VocabularyFilter{ID: []uint64{returned.ID}})
				require.NoError(t, err)
				require.Len(t, got, 1)
				assert.Equal(t, tt.data.Name, got[0].Data.Name)
				assert.Equal(t, tt.data.Title, got[0].Data.Title)
				assert.Equal(t, pointer.GetString(tt.data.Description), *got[0].Data.Description)
				assert.Equal(t, tt.data.ParentID, got[0].Data.ParentID)
			}
		})
	}
}

func TestVocabulary_Update(t *testing.T) {
	tests := []struct {
		name   string
		create func(t *testing.T, vocabulary repository.Vocabulary)
		update func(t *testing.T, vocabulary repository.Vocabulary)
		check  func(t *testing.T, client *ent.Client)
	}{
		{
			`ok`,
			func(t *testing.T, vocabulary repository.Vocabulary) {
				_, err := vocabulary.Create(context.TODO(), &model.VocabularyData{
					Name: `Jamaica`,
				})
				require.NoError(t, err)
			},
			func(t *testing.T, vocabulary repository.Vocabulary) {
				_, err := vocabulary.Update(context.TODO(), 1, &model.VocabularyData{
					Name: `Aruba`,
				})
				require.NoError(t, err)
			},
			func(t *testing.T, client *ent.Client) {
				row := client.Vocabulary.GetX(context.TODO(), 1)
				assert.Equal(t, `Aruba`, row.Name)
			},
		},
		{
			`set parentid nil`,
			func(t *testing.T, vocabulary repository.Vocabulary) {
				{
					_, err := vocabulary.Create(context.TODO(), &model.VocabularyData{
						Name: `Bermuda`,
					})
					require.NoError(t, err)
				}
				{
					_, err := vocabulary.Create(context.TODO(), &model.VocabularyData{
						Name:     `Bahama`,
						ParentID: pointer.ToUint64(1),
					})
					require.NoError(t, err)
				}
			},
			func(t *testing.T, vocabulary repository.Vocabulary) {
				_, err := vocabulary.Update(context.TODO(), 2, &model.VocabularyData{
					Name:     `Bahama`,
					ParentID: nil,
				})
				require.NoError(t, err)
			},
			func(t *testing.T, client *ent.Client) {
				row := client.Vocabulary.GetX(context.TODO(), 1)
				assert.Equal(t, pointer.ToUint64OrNil(0), row.ParentID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, client := vocabularyClient(t)
			tt.create(t, c)
			tt.update(t, c)
			tt.check(t, client)
		})
	}
}

func vocabularyClient(t *testing.T) (repository.Vocabulary, *ent.Client) {
	var client *ent.Client

	c := repo.NewVocabulary(func(t *testing.T) *ent.VocabularyClient {
		client = enttest.Open(t, "sqlite3", ":memory:?_fk=1", []enttest.Option{
			enttest.WithOptions(ent.Log(t.Log)),
		}...) // .Debug()

		t.Cleanup(func() {
			require.NoError(t, client.Close())
		})

		return client.Vocabulary
	}(t))

	return c, client
}
