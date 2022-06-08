package cmd

import (
	"github.com/spf13/cobra"
	"tagservice/cmd/loader"
)

func service(cmd *cobra.Command) *loader.Service {
	// Get DSN
	dsn, err := cmd.Flags().GetString(`dsn`)
	CheckErr(err)
	// Get verbose flag
	verbose, err := cmd.Flags().GetBool(`verbose`)
	CheckErr(err)
	// Load service
	service, err := loader.Load(cmd.Context(), dsn, verbose)
	CheckErr(err)
	return service
}

// https://github.com/spf13/cobra/pull/1568
func CheckErr(msg interface{}) {
	if msg != nil {
		panic(msg)
	}
}
