/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command

func serveCommand() *cobra.Command {
	var serveCmd = &cobra.Command{
		Use: "serve",
	}

	serveCmd.PersistentFlags().IntP("port", "p", 8080, "port on which the server will listen")
	CheckErr(serveCmd.MarkPersistentFlagRequired(`port`))

	serveCmd.AddCommand(&cobra.Command{
		Use:   "graphql",
		Short: "Run graphql server",
		Run: func(cmd *cobra.Command, args []string) {

		},
	})
	return serveCmd
}
