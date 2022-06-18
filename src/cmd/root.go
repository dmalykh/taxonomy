package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	// Root
	c := &cobra.Command{
		Use:   "tagservice",
		Short: "Service for tag management. It allows CRUD operations with tags, categories and namespaces.",
		Long:  `tagservice using for manage, build and debug your tags`,
	}
	defaultDSN := func() string {
		if dsn := os.Getenv(`DSN`); dsn != `` {
			return dsn
		}

		return `sqlite://./tagservice.db?cache=shared&_fk=1`
	}()

	c.PersistentFlags().String("dsn", defaultDSN, "Data source name (connection information)")
	c.PersistentFlags().BoolP("verbose", "v", false, "Make some output more verbose.")

	// Add subcommands
	c.AddCommand(initCommand(), categoryCommand(), tagCommand(), namespaceCommand(), relCommand(), serveCommand())

	return c
}

func Execute(c *cobra.Command) {
	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}
