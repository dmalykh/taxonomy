package repository_test

import (
	"context"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/dmalykh/tagservice/repository/entgo/ent"
	"github.com/dmalykh/tagservice/repository/entgo/ent/enttest"
	repo "github.com/dmalykh/tagservice/repository/entgo/repository"
	"github.com/dmalykh/tagservice/tagservice/model"
	"github.com/dmalykh/tagservice/tagservice/repository"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//goland:noinspection GoContextTodo,GoContextTodo,GoContextTodo,GoContextTodo,GoContextTodo,GoContextTodo
func TestCategory_Create(t *testing.T) {
	//goland:noinspection GoImportUsedAsName
	faker := faker.New()

	tests := []struct {
		name    string
		data    model.CategoryData
		prepare func(c *ent.Client)
		wantErr assert.ErrorAssertionFunc // returns continue
		check   func(t assert.TestingT, c *ent.Client)
	}{
		{
			`ok`,
			model.CategoryData{
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
			model.CategoryData{
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
			model.CategoryData{
				Title: faker.Beer().Name(),
			},
			func(c *ent.Client) {},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.ErrorIs(t, err, repository.ErrCreateCategory)

				return false
			},
			func(t assert.TestingT, c *ent.Client) {},
		},
		{
			`name duplicated, parentid unique, no error`,
			model.CategoryData{
				Name:  `gogolek`,
				Title: `kkk`,
				ParentID: func() *uint {
					var p uint = 2

					return &p
				}(),
			},
			func(c *ent.Client) {
				// https://www.sqlite.org/nulls.html
				c.Category.Create().SetName(`https://www.sqlite.org/nulls.html`).SetTitle(`NULL in SQL can't be UNIQ`).SaveX(context.TODO()) // id: 1
				c.Category.Create().SetName(`gogolek`).SetTitle(`kkk`).SetParentID(1).SaveX(context.TODO())                                  // id: 2
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.NoError(t, err)

				return true
			},
			func(t assert.TestingT, c *ent.Client) {
				// Additional check. It was a bug when parent was updated too
				category := c.Category.GetX(context.TODO(), 1)
				assert.Empty(t, category.ParentID)
			},
		},
		{
			`name and parentid duplicated, error`,
			model.CategoryData{
				Name:  `gogolek`,
				Title: `kkk`,
				ParentID: func() *uint {
					var p uint = 1

					return &p
				}(),
			},
			func(c *ent.Client) {
				// https://www.sqlite.org/nulls.html
				c.Category.Create().SetName(`https://www.sqlite.org/nulls.html`).SetTitle(`NULL in SQL can't be UNIQ`).SaveX(context.TODO())
				c.Category.Create().SetName(`gogolek`).SetTitle(`kkk`).SetParentID(1).SaveX(context.TODO())
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
			c, client := categoryClient(t)
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
				got, err := c.GetByID(ctx, returned.ID)
				require.NoError(t, err)
				assert.Equal(t, tt.data.Name, got.Data.Name)
				assert.Equal(t, tt.data.Title, got.Data.Title)
				assert.Equal(t, pointer.GetString(tt.data.Description), *got.Data.Description)
				assert.Equal(t, tt.data.ParentID, got.Data.ParentID)
			}
		})
	}
}

//goland:noinspection GoContextTodo,GoContextTodo,GoContextTodo,GoContextTodo,GoContextTodo,GoContextTodo,GoContextTodo
func TestCategory_Update(t *testing.T) {
	tests := []struct {
		name   string
		create func(t *testing.T, category repository.Category)
		update func(t *testing.T, category repository.Category)
		check  func(t *testing.T, client *ent.Client)
	}{
		{
			`ok`,
			func(t *testing.T, category repository.Category) {
				_, err := category.Create(context.TODO(), &model.CategoryData{
					Name: `Jamaica`,
				})
				require.NoError(t, err)
			},
			func(t *testing.T, category repository.Category) {
				_, err := category.Update(context.TODO(), 1, &model.CategoryData{
					Name: `Aruba`,
				})
				require.NoError(t, err)
			},
			func(t *testing.T, client *ent.Client) {
				row := client.Category.GetX(context.TODO(), 1)
				assert.Equal(t, `Aruba`, row.Name)
			},
		},
		{
			`set parentid nil`,
			func(t *testing.T, category repository.Category) {
				{
					_, err := category.Create(context.TODO(), &model.CategoryData{
						Name: `Bermuda`,
					})
					require.NoError(t, err)
				}
				{
					_, err := category.Create(context.TODO(), &model.CategoryData{
						Name:     `Bahama`,
						ParentID: pointer.ToUint(1),
					})
					require.NoError(t, err)
				}
			},
			func(t *testing.T, category repository.Category) {
				_, err := category.Update(context.TODO(), 2, &model.CategoryData{
					Name:     `Bahama`,
					ParentID: nil,
				})
				require.NoError(t, err)
			},
			func(t *testing.T, client *ent.Client) {
				row := client.Category.GetX(context.TODO(), 1)
				assert.Equal(t, pointer.ToIntOrNil(0), row.ParentID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, client := categoryClient(t)
			tt.create(t, c)
			tt.update(t, c)
			tt.check(t, client)
		})
	}
}

func categoryClient(t *testing.T) (repository.Category, *ent.Client) {
	var client *ent.Client

	c := repo.NewCategory(func(t *testing.T) *ent.CategoryClient {
		client = enttest.Open(t, "sqlite3", ":memory:?_fk=1", []enttest.Option{
			enttest.WithOptions(ent.Log(t.Log)),
		}...) // .Debug()

		t.Cleanup(func() {
			require.NoError(t, client.Close())
		})

		return client.Category
	}(t))

	return c, client
}
