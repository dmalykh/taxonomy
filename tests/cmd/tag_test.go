package cmd_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/dmalykh/taxonomy/internal/repository/entgo"
	ent2 "github.com/dmalykh/taxonomy/internal/repository/entgo/ent"
	"github.com/dmalykh/taxonomy/internal/repository/entgo/ent/enttest"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	cmd2 "github.com/dmalykh/taxonomy/cmd"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	suitetest "github.com/stretchr/testify/suite"
)

type TestTermOperations struct {
	suitetest.Suite
	dbpath string
	dsn    string
	client *ent2.Client
}

// nolint:gosec
func (suite *TestTermOperations) SetupTest() {
	rand.Seed(time.Now().UnixNano())
	suite.dbpath = suite.T().TempDir() + fmt.Sprintf(`cachedb%d.db`, rand.Int())
	suite.dsn = fmt.Sprintf(`sqlite://%s?mode=memory&cache=shared&_fk=1`, suite.dbpath)

	suite.client = enttest.Open(suite.T(), "sqlite3", fmt.Sprintf(`file:%s?_fk=1`, suite.dbpath), []enttest.Option{
		enttest.WithOptions(ent2.Log(suite.T().Log)),
	}...).Debug()

	cmd := cmd2.New()

	cmd.SetArgs([]string{`--dsn`, suite.dsn, `init`})
	suite.Require().NoError(cmd.Execute())
}

func (suite *TestTermOperations) TearDownTest() {
	if suite.dbpath != `` {

		os.Remove(suite.dbpath)
	}
}

func (suite *TestTermOperations) TestCreate() {
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
				suite.client.Vocabulary.Create().SetName(`test`).SaveX(context.TODO())
			},
			[][]string{{`create`, `Hello!`, `--vocabulary`, `1`}},
			func(err error) {
				suite.NoError(err)
			},
			func(out string) {
				suite.Empty(out)
			},
		},
		{
			`no vocabulary`,
			func() {},
			[][]string{{`create`, `Hello!`}},
			func(err error) {
				suite.Error(err)
				suite.Contains(err.Error(), `required flag(s) "vocabulary" not set`)
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
				cmd.SetArgs(append([]string{`--dsn`, suite.dsn, `term`}, command...))
				tt.Error(cmd.Execute())
			}
			out, _ := ioutil.ReadAll(b)
			tt.check(string(out))
		})
	}
}

func (suite *TestTermOperations) TestErrUpdate() {
	tests := []struct {
		name       string
		prepare    func()
		createArgs [][]string
		updateArgs [][]string
		check      func(t *TestTermOperations, out string)
	}{
		{
			`no error`,
			func() {
				suite.client.Vocabulary.Create().SetName(`test`).SaveX(context.TODO())
				suite.client.Vocabulary.Create().SetName(`test2`).SaveX(context.TODO())
			},
			[][]string{{`ohmyname`, `--vocabulary`, `1`}, {`cherrypie`, `--vocabulary`, `1`}},
			[][]string{{`1`, `--name`, `itsme`}, {`2`, `--vocabulary`, `2`}},
			func(t *TestTermOperations, out string) {
				c, err := entgo.Connect(context.TODO(), suite.dsn, false)
				cmd2.CheckErr(err)
				all := c.Term.Query().AllX(context.TODO())
				assert.Equal(t.T(), `itsme`, all[0].Name)
				assert.Equal(t.T(), 1, all[0].VocabularyID)
				assert.Equal(t.T(), `cherrypie`, all[1].Name)
				assert.Equal(t.T(), 2, all[1].VocabularyID)
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
					cmd.SetArgs(append([]string{`--dsn`, suite.dsn, `term`, `create`}, arg...))
					suite.Require().NoError(cmd.Execute())
					out, _ := ioutil.ReadAll(b)
					suite.Empty(out)
				}(arg)
			}

			var finalOutput []byte
			for _, arg := range tt.updateArgs {
				func(arg []string) {
					cmd := newcmd()
					cmd.SetArgs(append([]string{`--dsn`, suite.dsn, `term`, `update`}, arg...))
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

func TestTermOperationsSuite(t *testing.T) {
	suitetest.Run(t, new(TestTermOperations))
}
