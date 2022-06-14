/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS> @TODO

*/
package cmd

import (
	"github.com/dmalykh/tagservice/tagservice/model"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"strconv"
)

func tagCommand() *cobra.Command {

	var tagCmd = &cobra.Command{
		Use:   `tag`,
		Short: `CRUD operations with tags`,
		Run: func(cmd *cobra.Command, args []string) {
			CheckErr(cmd.Help())
		},
	}

	var createCmd = &cobra.Command{
		Use:        `create [name]`,
		Args:       cobra.ExactArgs(1),
		ArgAliases: []string{`name`},
		Short:      `Create new tag`,
		Run: func(cmd *cobra.Command, args []string) {
			_, err := service(cmd).Tag.Create(cmd.Context(), &model.TagData{
				Name:        args[0],
				Title:       cmd.Flag(`title`).Value.String(),
				Description: cmd.Flag(`description`).Value.String(),
				CategoryId: func() uint {
					if !cmd.Flags().Changed(`category`) {
						cobra.CompErrorln(`category required`)
					}
					categoryId, err := cmd.Flags().GetUint(`category`)
					CheckErr(err)
					return categoryId
				}(),
			})
			CheckErr(err)
		},
	}
	createCmd.Flags().StringP(`title`, `t`, ``, `title of the tag`)
	createCmd.Flags().Uint(`category`, 0, `id of category for the tag`)
	createCmd.Flags().String(`description`, ``, `description for the tag`)
	CheckErr(createCmd.MarkFlagRequired(`category`))

	var updateCmd = &cobra.Command{
		Use:   `update [id]`,
		Args:  cobra.ExactArgs(1),
		Short: `Update tag`,
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.ParseUint(args[0], 10, 32)
			CheckErr(err)
			var update model.TagData
			{
				name, err := cmd.Flags().GetString(`name`)
				CheckErr(err)
				if name != `` {
					update.Name = name
				}
			}
			{
				title, err := cmd.Flags().GetString(`title`)
				CheckErr(err)
				if title != `` {
					update.Title = title
				}
			}
			{
				description, err := cmd.Flags().GetString(`description`)
				CheckErr(err)
				if description != `` {
					update.Description = description
				}
			}
			{
				if cmd.Flags().Changed(`category`) {
					category, err := cmd.Flags().GetUint(`category`)
					CheckErr(err)
					if category > 0 {
						update.CategoryId = category
					}
				}
			}
			_, err = service(cmd).Tag.Update(cmd.Context(), uint(id), &update)
			CheckErr(err)
		},
	}
	updateCmd.Flags().Uint(`category`, 0, `id of category for this tag`)
	updateCmd.Flags().StringP(`name`, `n`, ``, `name of the tag`)
	updateCmd.Flags().StringP(`title`, `t`, ``, `title of the tag`)
	updateCmd.Flags().String(`description`, ``, `description for this category`)

	var deleteCmd = &cobra.Command{
		Use:   `delete [id]`,
		Args:  cobra.ExactArgs(1),
		Short: `Delete tag`,
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.ParseUint(args[0], 10, 32)
			CheckErr(err)
			CheckErr(service(cmd).Tag.Delete(cmd.Context(), uint(id)))
		},
	}

	var listCmd = &cobra.Command{
		Use:   `list [category's id] [limit] [offset]`,
		Args:  cobra.MinimumNArgs(1),
		Short: `Show all tags`,
		Run: func(cmd *cobra.Command, args []string) {
			categoryId, err := strconv.ParseUint(args[0], 10, 32)
			CheckErr(err)
			limit, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				limit = 10
			}
			offset, err := strconv.ParseUint(args[2], 10, 32)
			if err != nil {
				offset = 0
			}
			tags, err := service(cmd).Tag.GetList(cmd.Context(), uint(categoryId), uint(limit), uint(offset))
			CheckErr(err)

			table := tablewriter.NewWriter(cmd.OutOrStdout())
			table.SetHeader([]string{`ID`, `Name`, `Title`})

			for _, tag := range tags {
				table.Append(func(tag model.Tag) []string {
					return []string{
						strconv.Itoa(int(tag.Id)),
						tag.Data.Name,
						tag.Data.Title,
					}
				}(tag))
			}
			table.Render()
		},
	}

	tagCmd.AddCommand(createCmd, updateCmd, deleteCmd, listCmd)

	return tagCmd
}
