package repository

import (
	"context"
	"fmt"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"tagservice/repository/entgo/ent"
	"tagservice/repository/entgo/ent/enttest"
	"tagservice/server/model"
	"testing"
)

func TestCategory_Create(t *testing.T) {
	var faker = faker.New()
	tests := []struct {
		name    string
		data    model.CategoryData
		wantErr assert.ErrorAssertionFunc //returns continue
	}{
		{
			`ok`,
			model.CategoryData{
				Name:  faker.Beer().Name(),
				Title: faker.Beer().Name(),
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			`name empty, error`,
			model.CategoryData{
				Title: faker.Beer().Name(),
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.ErrorIs(t, err, ErrCreateCategory)
				return false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx = context.TODO()
			c := &Category{
				client: func(t *testing.T) *ent.CategoryClient {
					var client = enttest.Open(t, "sqlite3", ":memory:?_fk=1", []enttest.Option{
						enttest.WithOptions(ent.Log(t.Log)),
					}...).Debug()

					t.Cleanup(func() {
						require.NoError(t, client.Close())
					})
					return client.Category
				}(t),
			}
			returned, err := c.Create(ctx, &tt.data)
			if !tt.wantErr(t, err) {
				return
			}
			assert.EqualValues(t, tt.data, returned.Data)

			{
				got, err := c.GetById(ctx, returned.Id)
				require.NoError(t, err)
				assert.EqualValues(t, tt.data.Name, got.Data.Name)
				assert.EqualValues(t, tt.data.Title, got.Data.Title)
				assert.EqualValues(t, tt.data.Description, got.Data.Description)
				assert.EqualValues(t, tt.data.ParentId, got.Data.ParentId)
			}
		})
	}
}

func TestCategory_DeleteById(t *testing.T) {
	type fields struct {
		client *ent.CategoryClient
	}
	type args struct {
		ctx context.Context
		id  uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Category{
				client: tt.fields.client,
			}
			tt.wantErr(t, c.DeleteById(tt.args.ctx, tt.args.id), fmt.Sprintf("DeleteById(%v, %v)", tt.args.ctx, tt.args.id))
		})
	}
}
