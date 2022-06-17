//nolint:dupl
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

	"github.com/AlekSi/pointer"
	cmd2 "github.com/dmalykh/tagservice/cmd"
	"github.com/dmalykh/tagservice/repository/entgo"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	suitetest "github.com/stretchr/testify/suite"
)

type TestCategoryOperations struct {
	suitetest.Suite
	dbpath string
	dsn    string
}

//nolint:gosec
func (suite *TestCategoryOperations) SetupTest() {
	rand.Seed(time.Now().UnixNano())
	suite.dbpath = suite.T().TempDir() + fmt.Sprintf(`cachedb%d.db`, rand.Int())
	suite.dsn = fmt.Sprintf(`sqlite://%s?mode=memory&cache=shared&_fk=1`, suite.dbpath)

	cmd := cmd2.New()

	cmd.SetArgs([]string{`--dsn`, suite.dsn, `init`})
	suite.Require().NoError(cmd.Execute())
}

func (suite *TestCategoryOperations) TearDownTest() {
	if suite.dbpath != `` {
		//goland:noinspection GoUnhandledErrorResult
		os.Remove(suite.dbpath)
	}
}

func (suite *TestCategoryOperations) TestCreate() {
	var (
		cmd = cmd2.New()
		b   = bytes.NewBufferString(``)
	)

	cmd.SetOut(b)

	cmd.SetArgs([]string{`--dsn`, suite.dsn, `category`, `create`, `Hello!`})
	suite.NoError(cmd.Execute())

	out, _ := ioutil.ReadAll(b)
	suite.Empty(out)
}

func (suite *TestCategoryOperations) TestErrDuplicateNames() {
	tests := []struct {
		name       string
		firstArgs  [][]string
		secondArgs [][]string
		check      func(t *TestCategoryOperations, out string)
	}{
		{
			`duplicates error`,
			[][]string{{`ohmyname`}},
			[][]string{{`ohmyname`, `-p`, `1`}, {`ohmyname`, `-p`, `1`}},
			func(t *TestCategoryOperations, out string) {
				assert.Contains(t.T(), out, `must be unique`)
			},
		},
		{
			`no error with various names`,
			[][]string{{`ohmyname`}},
			[][]string{{`supername`}},
			func(t *TestCategoryOperations, out string) {
				assert.Empty(t.T(), out)
			},
		},
		{
			`no error with various parents`,
			[][]string{{`ohmyname`}},
			[][]string{{`ohmyname`, `--parent`, `1`}},
			func(t *TestCategoryOperations, out string) {
				assert.Empty(t.T(), out)
			},
		},
	}

	for _, tt := range tests {
		suite.TearDownTest()
		suite.SetupTest()
		suite.Run(tt.name, func() {
			cmd := cmd2.New()
			b := bytes.NewBufferString(``)
			cmd.SetOut(b)
			cmd.SetErr(b)

			// Catch panic in error
			defer func() {
				r := recover()
				if r != nil {
					tt.check(suite, fmt.Sprintf(`%v`, r))
				}
			}()

			for _, arg := range tt.firstArgs {
				cmd.SetArgs(append([]string{`--dsn`, suite.dsn, `category`, `create`}, arg...))
				suite.Require().NoError(cmd.Execute())
				out, _ := ioutil.ReadAll(b)
				suite.Empty(out)
			}

			var finalOutput []byte
			for _, arg := range tt.secondArgs {
				cmd.SetArgs(append([]string{`--dsn`, suite.dsn, `category`, `create`}, arg...))
				suite.Require().NoError(cmd.Execute())
				out, _ := ioutil.ReadAll(b)
				suite.Empty(out)
				finalOutput = out
			}
			tt.check(suite, string(finalOutput))
		})
	}
}

//goland:noinspection GoContextTodo,GoContextTodo
func (suite *TestCategoryOperations) TestErrUpdate() {
	tests := []struct {
		name       string
		createArgs [][]string
		updateArgs [][]string
		check      func(t *TestCategoryOperations, out string)
	}{
		{
			`no error`,
			[][]string{{`ohmyname`}, {`cherrypie`, `-p`, `1`}},
			[][]string{{`1`, `--name`, `itsme`}, {`2`, `--parent`, `-1`}, {`1`, `--parent`, `2`}},
			func(t *TestCategoryOperations, out string) {
				c, err := entgo.Connect(context.TODO(), suite.dsn, false)
				cmd2.CheckErr(err)
				all := c.Category.Query().AllX(context.TODO())
				assert.Equal(t.T(), `itsme`, all[0].Name)
				assert.Equal(t.T(), 2, *all[0].ParentID)
				assert.Equal(t.T(), `cherrypie`, all[1].Name)
				assert.Equal(t.T(), pointer.ToIntOrNil(0), all[1].ParentID)
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

			for _, arg := range tt.createArgs {
				func(arg []string) {
					cmd := newcmd()
					cmd.SetArgs(append([]string{`--dsn`, suite.dsn, `category`, `create`}, arg...))
					suite.Require().NoError(cmd.Execute())
					out, _ := ioutil.ReadAll(b)
					suite.Empty(out)
				}(arg)
			}

			var finalOutput []byte
			for _, arg := range tt.updateArgs {
				func(arg []string) {
					cmd := newcmd()
					cmd.SetArgs(append([]string{`--dsn`, suite.dsn, `category`, `update`}, arg...))
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

func TestCategoryOperationsSuite(t *testing.T) {
	suitetest.Run(t, new(TestCategoryOperations))
}
