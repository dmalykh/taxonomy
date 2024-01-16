package cmd

import (
	"github.com/dmalykh/taxonomy/taxonomy/model"
	"strconv"
	"unsafe"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func vocabularyCommand() *cobra.Command {
	vocabularyCmd := &cobra.Command{
		Use:   `vocabulary`,
		Short: `Operations with categories`,
		Run: func(cmd *cobra.Command, args []string) {
			CheckErr(cmd.Help())
		},
	}

	createCmd := &cobra.Command{
		Use:        `create [name]`,
		Args:       cobra.ExactArgs(1),
		ArgAliases: []string{`name`},
		Short:      `Create new vocabulary`,
		Long:       `Vocabulary's name and parent must be unique.`,
		Run: func(cmd *cobra.Command, args []string) {
			description := cmd.Flag(`description`).Value.String()
			_, err := service(cmd).Vocabulary.Create(cmd.Context(), &model.VocabularyData{
				Name:        args[0],
				Title:       cmd.Flag(`title`).Value.String(),
				Description: &description,
				ParentID: func() *uint64 {
					if !cmd.Flags().Changed(`parent`) {
						return nil
					}
					parentID, err := cmd.Flags().GetUint64(`parent`)
					CheckErr(err)
					id := parentID

					return &id
				}(),
			})
			CheckErr(err)
		},
	}

	createCmd.Flags().StringP(`title`, `t`, ``, `title of this vocabulary`)
	createCmd.Flags().UintP(`parent`, `p`, 0, `id of parent vocabulary for this vocabulary`)
	createCmd.Flags().String(`description`, ``, `description for this vocabulary`)

	updateCmd := &cobra.Command{
		Use:   `update [id]`,
		Args:  cobra.ExactArgs(1),
		Short: `Update vocabulary`,
		Long:  `Vocabulary's name and parent must be unique.`,
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.ParseUint(args[0], 10, 32)
			CheckErr(err)
			var update model.VocabularyData
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
						update.ParentID = (*uint64)(unsafe.Pointer(&parent))
					} else if parent < 0 {
						update.ParentID = func() *uint64 {
							var p uint64

							return &p
						}()
					}
				}
			}
			_, err = service(cmd).Vocabulary.Update(cmd.Context(), uint(id), &update)
			CheckErr(err)
		},
	}

	updateCmd.Flags().IntP(`parent`, `p`, 0, `id of parent vocabulary for this vocabulary (set negative value for nil value, i.e. -1)`)
	updateCmd.Flags().StringP(`name`, `n`, ``, `name of this vocabulary (name must be unique)`)
	updateCmd.Flags().StringP(`title`, `t`, ``, `title of this vocabulary`)
	updateCmd.Flags().String(`description`, ``, `description for this vocabulary`)

	deleteCmd := &cobra.Command{
		Use:   `delete [id]`,
		Args:  cobra.ExactArgs(1),
		Short: `Delete vocabulary`,
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.ParseUint(args[0], 10, 32)
			CheckErr(err)
			CheckErr(service(cmd).Vocabulary.Delete(cmd.Context(), uint(id)))
		},
	}

	listCmd := &cobra.Command{
		Use:   `list`,
		Args:  cobra.NoArgs,
		Short: `Show all vocabularies`,
		Run: func(cmd *cobra.Command, args []string) {
			vocabularies, err := service(cmd).Vocabulary.Get(cmd.Context(), &model.VocabularyFilter{})
			CheckErr(err)

			table := tablewriter.NewWriter(cmd.OutOrStdout())
			table.SetHeader([]string{`ID`, `Name`, `Title`, `Parent ID`})

			for _, vocabulary := range vocabularies {
				table.Append(func(vocabulary model.Vocabulary) []string {
					return []string{
						strconv.Itoa(int(vocabulary.ID)),
						vocabulary.Data.Name,
						vocabulary.Data.Title,
						func(parentId *uint64) string {
							if parentId == nil {
								return `—`
							}

							return strconv.Itoa(int(*parentId))
						}(vocabulary.Data.ParentID),
					}
				}(vocabulary))
			}
			table.Render()
		},
	}

	vocabularyCmd.AddCommand(createCmd, updateCmd, deleteCmd, listCmd)

	return vocabularyCmd
}
