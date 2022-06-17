package cmd

import (
	"strconv"

	"github.com/dmalykh/tagservice/tagservice/model"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func relCommand() *cobra.Command {
	relCmd := &cobra.Command{
		Use:   `rel`,
		Short: `Work with relations`,
		Run: func(cmd *cobra.Command, args []string) {
			CheckErr(cmd.Help())
		},
	}

	setCmd := &cobra.Command{
		Use:   `set`,
		Args:  cobra.NoArgs,
		Short: `Set relation`,
		Run: func(cmd *cobra.Command, args []string) {
			tagID, err := cmd.Flags().GetUint(`tag`)
			CheckErr(err)
			namespace, err := cmd.Flags().GetString(`namespace`)
			CheckErr(err)
			entitiesID, err := cmd.Flags().GetUintSlice(`entity`)
			CheckErr(err)
			CheckErr(service(cmd).Tag.SetRelation(cmd.Context(), tagID, namespace, entitiesID...))
		},
	}

	setCmd.Flags().UintP(`tag`, `t`, 0, `tag's id'`)
	setCmd.Flags().StringP(`namespace`, `n`, ``, `description for the tag`)
	setCmd.Flags().UintSliceP(`entity`, `e`, nil, `entity's id for relation`)
	CheckErr(setCmd.MarkFlagRequired(`tag`))
	CheckErr(setCmd.MarkFlagRequired(`namespace`))
	CheckErr(setCmd.MarkFlagRequired(`entity`))

	listCmd := &cobra.Command{
		Use:   `list`,
		Args:  cobra.NoArgs,
		Short: `List of relations`,
		Run: func(cmd *cobra.Command, args []string) {
			namespace, err := cmd.Flags().GetString(`namespace`)
			CheckErr(err)
			relations, err := service(cmd).Tag.GetRelations(cmd.Context(), &model.EntityFilter{Namespace: []string{namespace}})

			CheckErr(err)
			table := tablewriter.NewWriter(cmd.OutOrStdout())
			table.SetHeader([]string{`Tag ID`, `Entity ID`})

			for _, relation := range relations {
				table.Append(func(relation model.Relation) []string {
					return []string{
						strconv.Itoa(int(relation.TagID)),
						strconv.Itoa(int(relation.EntityID)),
					}
				}(relation))
			}
			table.Render()
		},
	}

	listCmd.Flags().StringP(`namespace`, `n`, ``, `description for the tag`)
	CheckErr(setCmd.MarkFlagRequired(`namespace`))

	relCmd.AddCommand(setCmd, listCmd)

	return relCmd
}
