package cmd

import (
	"math"
	"strconv"

	"github.com/dmalykh/tagservice/tagservice/model"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func namespaceCommand() *cobra.Command {
	namespaceCmd := &cobra.Command{
		Use:   "namespace",
		Short: "CRUD operations with namespaces",
		Run: func(cmd *cobra.Command, args []string) {
			CheckErr(cmd.Help())
		},
	}

	createCmd := &cobra.Command{
		Use:        "create [name]",
		Args:       cobra.ExactArgs(1),
		ArgAliases: []string{`name`},
		Short:      "Create new namespace (name must be unique)",
		Run: func(cmd *cobra.Command, args []string) {
			_, err := service(cmd).Namespace.Create(cmd.Context(), args[0])
			CheckErr(err)
		},
	}

	updateCmd := &cobra.Command{
		Use:   "update [id] [name]",
		Args:  cobra.ExactArgs(2),
		Short: "Update namespace (name must be unique)",
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.ParseUint(args[0], 10, 32)
			CheckErr(err)
			_, err = service(cmd).Namespace.Update(cmd.Context(), uint(id), args[1])
			CheckErr(err)
		},
	}

	deleteCmd := &cobra.Command{
		Use:   "delete [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Delete namespace",
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.ParseUint(args[0], 10, 32)
			CheckErr(err)
			CheckErr(service(cmd).Namespace.Delete(cmd.Context(), uint(id)))
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Args:  cobra.NoArgs,
		Short: "Show all namespaces",
		Run: func(cmd *cobra.Command, args []string) {
			namespaces, err := service(cmd).Namespace.GetList(cmd.Context(), math.MaxUint, 0)
			CheckErr(err)

			table := tablewriter.NewWriter(cmd.OutOrStdout())
			table.SetHeader([]string{"ID", "Name"})

			for _, namespace := range namespaces {
				table.Append(func(namespace model.Namespace) []string {
					return []string{
						strconv.Itoa(int(namespace.ID)),
						namespace.Name,
					}
				}(namespace))
			}
			table.Render()
		},
	}

	namespaceCmd.AddCommand(createCmd, updateCmd, deleteCmd, listCmd)

	return namespaceCmd
}
