package cmd

import (
	"github.com/dmalykh/taxonomy/taxonomy/model"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func relCommand() *cobra.Command {
	relCmd := &cobra.Command{
		Use:   `rel`,
		Short: `Work with references`,
		Run: func(cmd *cobra.Command, args []string) {
			CheckErr(cmd.Help())
		},
	}

	setCmd := &cobra.Command{
		Use:   `set`,
		Args:  cobra.NoArgs,
		Short: `Set reference`,
		Run: func(cmd *cobra.Command, args []string) {
			termID, err := cmd.Flags().GetUint(`term`)
			CheckErr(err)
			namespace, err := cmd.Flags().GetString(`namespace`)
			CheckErr(err)
			entitiesID, err := cmd.Flags().GetUintSlice(`entity`)
			CheckErr(err)
			CheckErr(service(cmd).Term.SetReference(cmd.Context(), termID, namespace, entitiesID...))
		},
	}

	setCmd.Flags().UintP(`term`, `t`, 0, `term's id'`)
	setCmd.Flags().StringP(`namespace`, `n`, ``, `description for the term`)
	setCmd.Flags().UintSliceP(`entity`, `e`, nil, `entity's id for reference`)
	CheckErr(setCmd.MarkFlagRequired(`term`))
	CheckErr(setCmd.MarkFlagRequired(`namespace`))
	CheckErr(setCmd.MarkFlagRequired(`entity`))

	listCmd := &cobra.Command{
		Use:   `list`,
		Args:  cobra.NoArgs,
		Short: `List of references`,
		Run: func(cmd *cobra.Command, args []string) {
			namespace, err := cmd.Flags().GetString(`namespace`)
			CheckErr(err)
			references, err := service(cmd).Term.GetReferences(cmd.Context(), &model.EntityFilter{Namespace: []string{namespace}})

			CheckErr(err)
			table := tablewriter.NewWriter(cmd.OutOrStdout())
			table.SetHeader([]string{`Term ID`, `Entity ID`})

			for _, reference := range references {
				table.Append(func(reference model.Reference) []string {
					return []string{
						strconv.Itoa(int(reference.TermID)),
						reference.EntityID,
					}
				}(reference))
			}
			table.Render()
		},
	}

	listCmd.Flags().StringP(`namespace`, `n`, ``, `description for the term`)
	CheckErr(setCmd.MarkFlagRequired(`namespace`))

	relCmd.AddCommand(setCmd, listCmd)

	return relCmd
}
