/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS> @TODO

*/
package cmd

import (
	"github.com/dmalykh/tagservice/tagservice/model"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"strconv"
	"unsafe"
)

func categoryCommand() *cobra.Command {

	var categoryCmd = &cobra.Command{
		Use:   `category`,
		Short: `CRUD operations with categories`,
		Run: func(cmd *cobra.Command, args []string) {
			CheckErr(cmd.Help())
		},
	}

	var createCmd = &cobra.Command{
		Use:        `create [name]`,
		Args:       cobra.ExactArgs(1),
		ArgAliases: []string{`name`},
		Short:      `Create new category`,
		Long:       `Category's name and parent must be unique.`,
		Run: func(cmd *cobra.Command, args []string) {
			var description = cmd.Flag(`description`).Value.String()
			_, err := service(cmd).Category.Create(cmd.Context(), &model.CategoryData{
				Name:        args[0],
				Title:       cmd.Flag(`title`).Value.String(),
				Description: &description,
				ParentId: func() *uint {
					if !cmd.Flags().Changed(`parent`) {
						return nil
					}
					parentId, err := cmd.Flags().GetUint(`parent`)
					CheckErr(err)
					var id = uint(parentId)
					return &id
				}(),
			})
			CheckErr(err)
		},
	}
	createCmd.Flags().StringP(`title`, `t`, ``, `title of this category`)
	createCmd.Flags().UintP(`parent`, `p`, 0, `id of parent category for this category`)
	createCmd.Flags().String(`description`, ``, `description for this category`)

	var updateCmd = &cobra.Command{
		Use:   `update [id]`,
		Args:  cobra.ExactArgs(1),
		Short: `Update category`,
		Long:  `Category's name and parent must be unique.`,
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.ParseUint(args[0], 10, 32)
			CheckErr(err)
			var update model.CategoryData
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
				if cmd.Flags().Changed(`parent`) {
					parent, err := cmd.Flags().GetInt(`parent`)
					CheckErr(err)
					if parent > 0 {
						update.ParentId = (*uint)(unsafe.Pointer(&parent))
					} else if parent < 0 {
						update.ParentId = func() *uint {
							var p uint = 0
							return &p
						}()
					}
				}
			}
			_, err = service(cmd).Category.Update(cmd.Context(), uint(id), &update)
			CheckErr(err)
		},
	}
	updateCmd.Flags().IntP(`parent`, `p`, 0, `id of parent category for this category (set negative value for nil value, i.e. -1)`)
	updateCmd.Flags().StringP(`name`, `n`, ``, `name of this category (name must be unique)`)
	updateCmd.Flags().StringP(`title`, `t`, ``, `title of this category`)
	updateCmd.Flags().String(`description`, ``, `description for this category`)

	var deleteCmd = &cobra.Command{
		Use:   `delete [id]`,
		Args:  cobra.ExactArgs(1),
		Short: `Delete category`,
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.ParseUint(args[0], 10, 32)
			CheckErr(err)
			CheckErr(service(cmd).Category.Delete(cmd.Context(), uint(id)))
		},
	}

	var listCmd = &cobra.Command{
		Use:   `list`,
		Args:  cobra.NoArgs,
		Short: `Show all categorys`,
		Run: func(cmd *cobra.Command, args []string) {
			categorys, err := service(cmd).Category.GetList(cmd.Context(), &model.CategoryFilter{})
			CheckErr(err)

			table := tablewriter.NewWriter(cmd.OutOrStdout())
			table.SetHeader([]string{`ID`, `Name`, `Title`, `Parent ID`})

			for _, category := range categorys {
				table.Append(func(category model.Category) []string {
					return []string{
						strconv.Itoa(int(category.Id)),
						category.Data.Name,
						category.Data.Title,
						func(parentId *uint) string {
							if parentId == nil {
								return `—`
							}
							return strconv.Itoa(int(*parentId))
						}(category.Data.ParentId),
					}
				}(category))
			}
			table.Render()
		},
	}

	categoryCmd.AddCommand(createCmd, updateCmd, deleteCmd, listCmd)

	return categoryCmd
}
