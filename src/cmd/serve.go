/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/dmalykh/tagservice/api/graphql"
	"github.com/spf13/cobra"
	"strconv"
)

// serveCmd represents the serve command

func serveCommand() *cobra.Command {
	var serveCmd = &cobra.Command{
		Use:   `serve`,
		Short: `Run API server`,
	}

	serveCmd.PersistentFlags().IntP(`port`, `p`, 8080, `port on which the github.com/dmalykh/tagservice will listen`)
	CheckErr(serveCmd.MarkPersistentFlagRequired(`port`))

	serveCmd.AddCommand(&cobra.Command{
		Use:   `graphql`,
		Short: `Run graphql API`,
		Run: func(cmd *cobra.Command, args []string) {
			// Get port flag
			port, err := cmd.Flags().GetInt(`port`)
			CheckErr(err)
			// Get verbose flag
			verbose, err := cmd.Flags().GetBool(`verbose`)
			CheckErr(err)
			// Run service
			var s = service(cmd)
			CheckErr(graphql.Serve(&graphql.Config{
				Port:             strconv.Itoa(port),
				TagService:       s.Tag,
				CategoryService:  s.Category,
				NamespaceService: s.Namespace,
				Verbose:          verbose,
			}))
		},
	})
	return serveCmd
}
