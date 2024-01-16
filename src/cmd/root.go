package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	// Root
	c := &cobra.Command{
		Use:   "taxonomy",
		Short: "the service is used for manage taxonomy",
		Long:  `Service for taxonomy management. It allows CRUD operations with terms, vocabularies and namespaces.`,
	}
	defaultDSN := func() string {
		if dsn := os.Getenv(`DSN`); dsn != `` {
			return dsn
		}

		return `sqlite://./internal.db?cache=shared&_fk=1`
	}()

	c.PersistentFlags().String("dsn", defaultDSN, "Data source name (connection information)")
	c.PersistentFlags().BoolP("verbose", "v", false, "Make some output more verbose.")

	// Add subcommands
	c.AddCommand(initCommand(), vocabularyCommand(), termCommand(), namespaceCommand(), relCommand(), serveCommand())

	return c
}

func Execute(c *cobra.Command) {
	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}
