/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"tagservice/repository/entgo"
)

func initCommand() *cobra.Command {
	return &cobra.Command{
		Use:   `init`,
		Short: `Initiate service, create table in database`,
		Run: func(cmd *cobra.Command, args []string) {
			// Get DSN
			dsn, err := cmd.Flags().GetString(`dsn`)
			CheckErr(err)
			// Verbose
			v, err := cmd.Flags().GetBool(`verbose`)
			CheckErr(err)
			// Connect
			CheckErr(entgo.Init(cmd.Context(), dsn, v))
		},
	}
}
