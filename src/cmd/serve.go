package cmd

import (
	"strconv"

	"github.com/dmalykh/taxonomy/api/graphql"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command.
func serveCommand() *cobra.Command {
	serveCmd := &cobra.Command{
		Use:   `serve`,
		Short: `Run API server`,
	}

	serveCmd.PersistentFlags().IntP(`port`, `p`, 8080, `port on which the github.com/dmalykh/internal will listen`) //nolint:gomnd

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
			s := service(cmd)
			CheckErr(graphql.Serve(&graphql.Config{
				Port:              strconv.Itoa(port),
				TermService:       s.Term,
				VocabularyService: s.Vocabulary,
				NamespaceService:  s.Namespace,
				Verbose:           verbose,
			}))
		},
	})

	return serveCmd
}
