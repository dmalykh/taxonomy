package integration_tests

import (
	"bytes"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	suitetest "github.com/stretchr/testify/suite"
	"io/ioutil"
	"math/rand"
	"os"
	cmd2 "tagservice/cmd"
	"tagservice/repository/entgo"
	"testing"
	"time"
)

type TestNamespaceOperations struct {
	suitetest.Suite
	dbpath string
	dsn    string
}

func (suite *TestNamespaceOperations) SetupTest() {
	rand.Seed(time.Now().UnixNano())
	suite.dbpath = suite.T().TempDir() + fmt.Sprintf(`cachedb%d.db`, rand.Int())
	suite.dsn = fmt.Sprintf(`sqlite://%s?mode=memory&cache=shared&_fk=1`, suite.dbpath)

	var cmd = cmd2.New()
	cmd.SetArgs([]string{`--dsn`, suite.dsn, `init`})
	suite.Require().NoError(cmd.Execute())
}

func (suite *TestNamespaceOperations) TearDownTest() {
	if suite.dbpath != `` {
		os.Remove(suite.dbpath)
	}
}

func (suite *TestNamespaceOperations) TestCreate() {
	var cmd = cmd2.New()
	var b = bytes.NewBufferString(``)
	cmd.SetOut(b)

	cmd.SetArgs([]string{`--dsn`, suite.dsn, `namespace`, `create`, `Hello!`})
	suite.NoError(cmd.Execute())
	out, _ := ioutil.ReadAll(b)
	suite.Empty(out)
	c, err := entgo.Connect(context.TODO(), suite.dsn, false)
	cmd2.CheckErr(err)
	ns := c.Namespace.GetX(context.TODO(), 1)
	assert.Equal(suite.T(), `Hello!`, ns.Name)
}

func (suite *TestNamespaceOperations) TestUpdate() {
	var tests = []struct {
		name       string
		createArgs [][]string
		updateArgs [][]string
		check      func(t *TestNamespaceOperations, out string)
	}{
		{
			`no error`,
			[][]string{{`ohmyname`}},
			[][]string{{`1`, `aruba`}},
			func(t *TestNamespaceOperations, out string) {
				c, err := entgo.Connect(context.TODO(), suite.dsn, false)
				cmd2.CheckErr(err)
				var all = c.Namespace.Query().AllX(context.TODO())
				assert.Len(t.T(), all, 1)
				assert.Equal(t.T(), `aruba`, all[0].Name)
			},
		},
	}

	var b = bytes.NewBufferString(``)
	var newcmd = func() *cobra.Command {
		var c = cmd2.New()
		c.SetOut(b)
		c.SetErr(b)
		return c
	}

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
					var cmd = newcmd()
					cmd.SetArgs(append([]string{`--dsn`, suite.dsn, `namespace`, `create`}, arg...))
					suite.Require().NoError(cmd.Execute())
					out, _ := ioutil.ReadAll(b)
					suite.Empty(out)
				}(arg)
			}

			var finalOutput []byte
			for _, arg := range tt.updateArgs {
				func(arg []string) {
					var cmd = newcmd()
					cmd.SetArgs(append([]string{`--dsn`, suite.dsn, `namespace`, `update`}, arg...))
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

func TestNamespaceOperationsSuite(t *testing.T) {
	suitetest.Run(t, new(TestNamespaceOperations))
}
