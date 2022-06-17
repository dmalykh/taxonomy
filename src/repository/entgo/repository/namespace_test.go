package repository_test

import (
	"context"
	"testing"

	"github.com/dmalykh/tagservice/repository/entgo/ent"
	"github.com/dmalykh/tagservice/repository/entgo/ent/enttest"
	repo "github.com/dmalykh/tagservice/repository/entgo/repository"
	"github.com/dmalykh/tagservice/tagservice/repository"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//goland:noinspection GoContextTodo
func TestNamespace_Create(t *testing.T) {
	//goland:noinspection GoImportUsedAsName
	faker := faker.New()
	tests := []struct {
		name          string
		namespaceName string
		wantErr       assert.ErrorAssertionFunc // returns continue
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
			ctx := context.TODO()
			c := repo.NewNamespace(func(t *testing.T) *ent.NamespaceClient {
				client := enttest.Open(t, "sqlite3", ":memory:?_fk=1", []enttest.Option{
					enttest.WithOptions(ent.Log(t.Log)),
				}...).Debug()

				t.Cleanup(func() {
					require.NoError(t, client.Close())
				})

				return client.Namespace
			}(t))
			returned, err := c.Create(ctx, tt.namespaceName)
			if !tt.wantErr(t, err) {
				return
			}
			assert.EqualValues(t, tt.namespaceName, returned.Name)

			{
				got, err := c.GetByID(ctx, returned.ID)
				require.NoError(t, err)
				assert.EqualValues(t, tt.namespaceName, got.Name)
			}
		})
	}
}
