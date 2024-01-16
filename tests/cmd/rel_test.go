package cmd_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/dmalykh/taxonomy/internal/repository/entgo/ent"
	"github.com/dmalykh/taxonomy/internal/repository/entgo/ent/enttest"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	command "github.com/dmalykh/taxonomy/cmd"
	"github.com/jaswdr/faker"
	suitetest "github.com/stretchr/testify/suite"
)

type TestRelOperations struct {
	suitetest.Suite
	dbpath string
	dsn    string
	client *ent.Client
}

// nolint:gosec
func (suite *TestRelOperations) SetupTest() {
	rand.Seed(time.Now().UnixNano())
	suite.dbpath = suite.T().TempDir() + fmt.Sprintf(`cachedb%d.db`, rand.Int())
	suite.dsn = fmt.Sprintf(`sqlite://%s?mode=memory&cache=shared&_fk=1`, suite.dbpath)

	suite.client = enttest.Open(suite.T(), "sqlite3", fmt.Sprintf(`file:%s?_fk=1`, suite.dbpath), []enttest.Option{
		enttest.WithOptions(ent.Log(suite.T().Log)),
	}...).Debug()
}

func (suite *TestRelOperations) TearDownTest() {
	if suite.dbpath != `` {

		os.Remove(suite.dbpath)
	}
}

func (suite *TestRelOperations) TestSet() {

	faker := faker.New()

	tests := []struct {
		name     string
		prepare  func()
		commands [][]string
		check    func(out string)
	}{
		{
			`ok`,
			func() {
				vocabulary := suite.client.Vocabulary.Create().SetName(faker.Beer().Name()).SaveX(context.TODO())
				suite.client.Term.Create().SetName(faker.Beer().Name()).SetVocabularyID(vocabulary.ID).SaveX(context.TODO())
				suite.client.Namespace.Create().SetName(`gorilka`).SaveX(context.TODO())
			},
			[][]string{
				{`set`, `--term`, `1`, `--namespace`, `gorilka`, `--entity`, `333`, `--entity`, `444`, `--entity`, `333`},
				{`set`, `--term`, `1`, `--namespace`, `jenever`, `--entity`, `555`, `--entity`, `333`},
			},
			func(out string) {
				suite.Empty(out)
				all := suite.client.Reference.Query().AllX(context.TODO())
				suite.Len(all, 4)
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			cmd := command.New()
			tt.prepare()
			b := bytes.NewBufferString(``)
			cmd.SetOut(b)

			// Catch panic in error
			defer func() {
				r := recover()
				if r != nil {
					tt.check(fmt.Sprintf(`%v`, r))
				}
			}()

			for _, command := range tt.commands {
				cmd.SetArgs(append([]string{`--dsn`, suite.dsn, `rel`}, command...))
				suite.NoError(cmd.Execute())
			}
			out, _ := ioutil.ReadAll(b)
			tt.check(string(out))
		})
	}
}

func TestRelOperationsSuite(t *testing.T) {
	suitetest.Run(t, new(TestRelOperations))
}
