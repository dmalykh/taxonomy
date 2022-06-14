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
		Use:   "tagservice",
		Short: "Service for tag management. It allows CRUD operations with tags, categories and namespaces.",
		Long:  `tagservice is the cli to manage, build and debug your tags`,
	}
	c.PersistentFlags().String("dsn", os.Getenv(`DSN`), "Data source name (connection information)")
	c.PersistentFlags().BoolP("verbose", "v", false, "Make some output more verbose.")
	CheckErr(c.MarkPersistentFlagRequired("dsn"))

	// Add subcommands
	c.AddCommand(initCommand(), categoryCommand(), tagCommand(), namespaceCommand(), relCommand(), serveCommand())

	return c
}

func Execute(c *cobra.Command) {
	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}
