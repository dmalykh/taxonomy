package repository_test

import (
	"context"
	"github.com/dmalykh/taxonomy/internal/repository/entgo/ent"
	"github.com/dmalykh/taxonomy/internal/repository/entgo/ent/enttest"
	repo "github.com/dmalykh/taxonomy/internal/repository/entgo/repository"
	"github.com/dmalykh/taxonomy/taxonomy/model"
	"github.com/dmalykh/taxonomy/taxonomy/repository"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/xiaoqidun/entps"
	"testing"
)

func TestNamespace_Create(t *testing.T) {

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

			returned, err := c.Create(ctx, &model.NamespaceData{Name: tt.namespaceName})
			if !tt.wantErr(t, err) {
				return
			}
			assert.EqualValues(t, tt.namespaceName, returned.Data.Name)

			{
				got, err := c.Get(ctx, &repository.NamespaceFilter{ID: []uint64{returned.ID}})
				require.NoError(t, err)
				require.Len(t, got, 1)
				assert.EqualValues(t, tt.namespaceName, got[0].Data.Name)
			}
		})
	}
}
