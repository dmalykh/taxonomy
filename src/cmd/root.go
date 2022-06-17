/*
Copyright © 2022 Daniil Malykh daniil.malykh@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	// Root
	var c = &cobra.Command{
		Use:   "github.com/dmalykh/tagservice",
		Short: "Service for tag management. It allows CRUD operations with tags, categories and namespaces.",
		Long:  `github.com/dmalykh/tagservice is the cli to manage, build and debug your tags`,
	}
	var defaultDSN = func() string {
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
