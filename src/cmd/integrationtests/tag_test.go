package integrationtests_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	cmd2 "github.com/dmalykh/tagservice/cmd"
	"github.com/dmalykh/tagservice/repository/entgo"
	"github.com/dmalykh/tagservice/repository/entgo/ent"
	"github.com/dmalykh/tagservice/repository/entgo/ent/enttest"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	suitetest "github.com/stretchr/testify/suite"
)

type TestTagOperations struct {
	suitetest.Suite
	dbpath string
	dsn    string
	client *ent.Client
}

// nolint:gosec
func (suite *TestTagOperations) SetupTest() {
	rand.Seed(time.Now().UnixNano())
	suite.dbpath = suite.T().TempDir() + fmt.Sprintf(`cachedb%d.db`, rand.Int())
	suite.dsn = fmt.Sprintf(`sqlite://%s?mode=memory&cache=shared&_fk=1`, suite.dbpath)

	suite.client = enttest.Open(suite.T(), "sqlite3", fmt.Sprintf(`file:%s?_fk=1`, suite.dbpath), []enttest.Option{
		enttest.WithOptions(ent.Log(suite.T().Log)),
	}...).Debug()

	cmd := cmd2.New()

	cmd.SetArgs([]string{`--dsn`, suite.dsn, `init`})
	suite.Require().NoError(cmd.Execute())
}

func (suite *TestTagOperations) TearDownTest() {
	if suite.dbpath != `` {
		//goland:noinspection GoUnhandledErrorResult
		os.Remove(suite.dbpath)
	}
}

//goland:noinspection GoContextTodo
func (suite *TestTagOperations) TestCreate() {
	tests := []struct {
		name     string
		prepare  func()
		commands [][]string
		Error    func(err error)
		check    func(out string)
	}{
		{
			`ok`,
			func() {
				suite.client.Category.Create().SetName(`test`).SaveX(context.TODO())
			},
			[][]string{{`create`, `Hello!`, `--category`, `1`}},
			func(err error) {
				suite.NoError(err)
			},
			func(out string) {
				suite.Empty(out)
			},
		},
		{
			`no category`,
			func() {},
			[][]string{{`create`, `Hello!`}},
			func(err error) {
				suite.Error(err)
				suite.Contains(err.Error(), `required flag(s) "category" not set`)
			},
			func(out string) {
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			cmd := cmd2.New()
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
				cmd.SetArgs(append([]string{`--dsn`, suite.dsn, `tag`}, command...))
				tt.Error(cmd.Execute())
			}
			out, _ := ioutil.ReadAll(b)
			tt.check(string(out))
		})
	}
}

//goland:noinspection GoContextTodo,GoContextTodo,GoContextTodo,GoContextTodo
func (suite *TestTagOperations) TestErrUpdate() {
	tests := []struct {
		name       string
		prepare    func()
		createArgs [][]string
		updateArgs [][]string
		check      func(t *TestTagOperations, out string)
	}{
		{
			`no error`,
			func() {
				suite.client.Category.Create().SetName(`test`).SaveX(context.TODO())
				suite.client.Category.Create().SetName(`test2`).SaveX(context.TODO())
			},
			[][]string{{`ohmyname`, `--category`, `1`}, {`cherrypie`, `--category`, `1`}},
			[][]string{{`1`, `--name`, `itsme`}, {`2`, `--category`, `2`}},
			func(t *TestTagOperations, out string) {
				c, err := entgo.Connect(context.TODO(), suite.dsn, false)
				cmd2.CheckErr(err)
				all := c.Tag.Query().AllX(context.TODO())
				assert.Equal(t.T(), `itsme`, all[0].Name)
				assert.Equal(t.T(), 1, all[0].CategoryID)
				assert.Equal(t.T(), `cherrypie`, all[1].Name)
				assert.Equal(t.T(), 2, all[1].CategoryID)
			},
		},
	}

	var (
		b      = bytes.NewBufferString(``)
		newcmd = func() *cobra.Command {
			c := cmd2.New()
			c.SetOut(b)
			c.SetErr(b)

			return c
		}
	)

	for _, tt := range tests {
		suite.TearDownTest()
		suite.SetupTest()
		suite.Run(tt.name, func() {
			// Catch panic in error
			defer func() {
				r := recover()
				if r != nil {
					tt.check(suite, fmt.Sprintf(`%v`, r))
				}
				b.Reset()
			}()
			tt.prepare()
			for _, arg := range tt.createArgs {
				func(arg []string) {
					cmd := newcmd()
					cmd.SetArgs(append([]string{`--dsn`, suite.dsn, `tag`, `create`}, arg...))
					suite.Require().NoError(cmd.Execute())
					out, _ := ioutil.ReadAll(b)
					suite.Empty(out)
				}(arg)
			}

			var finalOutput []byte
			for _, arg := range tt.updateArgs {
				func(arg []string) {
					cmd := newcmd()
					cmd.SetArgs(append([]string{`--dsn`, suite.dsn, `tag`, `update`}, arg...))
					suite.Require().NoError(cmd.Execute())
					out, _ := ioutil.ReadAll(b)
					suite.Empty(out)
					finalOutput = out
				}(arg)
			}
			tt.check(suite, string(finalOutput))
		})
	}
}

func TestTagOperationsSuite(t *testing.T) {
	suitetest.Run(t, new(TestTagOperations))
}
