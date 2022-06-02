package repository

import (
	"context"
	"fmt"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"tagservice/repository/entgo/ent"
	"tagservice/repository/entgo/ent/enttest"
	"tagservice/server/repository"
	"testing"
)

func TestNamespace_Create(t *testing.T) {
	var faker = faker.New()
	tests := []struct {
		name          string
		namespaceName string
		wantErr       assert.ErrorAssertionFunc //returns continue
	}{
		{
			`ok`,
			faker.Beer().Name(),
			func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			`name empty, error`,
			``,
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.ErrorIs(t, err, repository.ErrCreateNamespace)
				return false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx = context.TODO()
			c := &Namespace{
				client: func(t *testing.T) *ent.NamespaceClient {
					var client = enttest.Open(t, "sqlite3", ":memory:?_fk=1", []enttest.Option{
						enttest.WithOptions(ent.Log(t.Log)),
					}...).Debug()

					t.Cleanup(func() {
						require.NoError(t, client.Close())
					})
					return client.Namespace
				}(t),
			}
			returned, err := c.Create(ctx, tt.namespaceName)
			if !tt.wantErr(t, err) {
				return
			}
			assert.EqualValues(t, tt.namespaceName, returned.Name)

			{
				got, err := c.GetById(ctx, returned.Id)
				require.NoError(t, err)
				assert.EqualValues(t, tt.namespaceName, got.Name)
			}
		})
	}
}

func TestNamespace_DeleteById(t *testing.T) {
	type fields struct {
		client *ent.NamespaceClient
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
			c := &Namespace{
				client: tt.fields.client,
			}
			tt.wantErr(t, c.DeleteById(tt.args.ctx, tt.args.id), fmt.Sprintf("DeleteById(%v, %v)", tt.args.ctx, tt.args.id))
		})
	}
}
