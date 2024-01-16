package cmd

import (
	"github.com/dmalykh/taxonomy/cmd/loader"
	"github.com/spf13/cobra"
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

// CheckErr check error and panics if error exists  https://github.com/spf13/cobra/pull/1568
func CheckErr(msg interface{}) {
	if msg != nil {
		panic(msg)
	}
}
