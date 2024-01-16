package cmd

import (
	"github.com/dmalykh/taxonomy/internal/repository/entgo"
	"github.com/spf13/cobra"
)

func initCommand() *cobra.Command {
	return &cobra.Command{
		Use:   `init`,
		Short: `Initiate service, create tables in a database`,
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
