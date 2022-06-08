package cmd

import (
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"strconv"
	"tagservice/server/model"
)

func relCommand() *cobra.Command {
	var relCmd = &cobra.Command{
		Use:   `rel`,
		Short: `Work with relations`,
		Run: func(cmd *cobra.Command, args []string) {
			CheckErr(cmd.Help())
		},
	}

	var setCmd = &cobra.Command{
		Use:   `set`,
		Args:  cobra.NoArgs,
		Short: `Set relation`,
		Run: func(cmd *cobra.Command, args []string) {
			tagId, err := cmd.Flags().GetUint(`tag`)
			CheckErr(err)
			namespace, err := cmd.Flags().GetString(`namespace`)
			CheckErr(err)
			entitiesId, err := cmd.Flags().GetUintSlice(`entity`)
			CheckErr(err)
			CheckErr(service(cmd).Tag.SetRelation(cmd.Context(), tagId, namespace, entitiesId...))
		},
	}
	setCmd.Flags().UintP(`tag`, `t`, 0, `tag's id'`)
	setCmd.Flags().StringP(`namespace`, `n`, ``, `description for the tag`)
	setCmd.Flags().UintSliceP(`entity`, `e`, nil, `entity's id for relation`)
	CheckErr(setCmd.MarkFlagRequired(`tag`))
	CheckErr(setCmd.MarkFlagRequired(`namespace`))
	CheckErr(setCmd.MarkFlagRequired(`entity`))

	var listCmd = &cobra.Command{
		Use:   `list`,
		Args:  cobra.NoArgs,
		Short: `List of relations`,
		Run: func(cmd *cobra.Command, args []string) {
			namespace, err := cmd.Flags().GetString(`namespace`)
			CheckErr(err)
			relations, err := service(cmd).Tag.GetRelationEntities(cmd.Context(), namespace, nil)

			CheckErr(err)
			table := tablewriter.NewWriter(cmd.OutOrStdout())
			table.SetHeader([]string{`Tag ID`, `Entity ID`})

			for _, relation := range relations {
				table.Append(func(relation model.Relation) []string {
					return []string{
						strconv.Itoa(int(relation.TagId)),
						strconv.Itoa(int(relation.EntityId)),
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
