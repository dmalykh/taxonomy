package cmd

import (
	"github.com/dmalykh/taxonomy/taxonomy/model"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func termCommand() *cobra.Command {
	termCmd := &cobra.Command{
		Use:   `term`,
		Short: `CRUD operations with terms`,
		Run: func(cmd *cobra.Command, args []string) {
			CheckErr(cmd.Help())
		},
	}

	createCmd := &cobra.Command{
		Use:        `create [name]`,
		Args:       cobra.ExactArgs(1),
		ArgAliases: []string{`name`},
		Short:      `Create new term`,
		Run: func(cmd *cobra.Command, args []string) {
			_, err := service(cmd).Term.Create(cmd.Context(), &model.TermData{
				Name:        args[0],
				Title:       cmd.Flag(`title`).Value.String(),
				Description: cmd.Flag(`description`).Value.String(),
				VocabularyID: func() []uint64 {
					if !cmd.Flags().Changed(`vocabulary`) {
						cobra.CompErrorln(`vocabulary required`)
					}
					vocabularyID, err := cmd.Flags().GetUintSlice(`vocabulary`)
					CheckErr(err)

					return vocabularyID
				}(),
			})
			CheckErr(err)
		},
	}

	createCmd.Flags().StringP(`title`, `t`, ``, `title of the term`)
	createCmd.Flags().UintSliceP(`vocabulary`, `v`, {}, `id of vocabulary for the term`)
	createCmd.Flags().String(`description`, ``, `description for the term`)
	CheckErr(createCmd.MarkFlagRequired(`vocabulary`))

	updateCmd := &cobra.Command{
		Use:   `update [id]`,
		Args:  cobra.ExactArgs(1),
		Short: `Update term`,
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.ParseUint(args[0], 10, 32)
			CheckErr(err)
			var update model.TermData
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
				if cmd.Flags().Changed(`vocabulary`) {
					vocabulary, err := cmd.Flags().GetUint64(`vocabulary`)
					CheckErr(err)
					if vocabulary > 0 {
						update.VocabularyID = vocabulary
					}
				}
			}
			_, err = service(cmd).Term.Update(cmd.Context(), uint(id), &update)
			CheckErr(err)
		},
	}

	updateCmd.Flags().Uint(`vocabulary`, 0, `id of vocabulary for this term`)
	updateCmd.Flags().StringP(`name`, `n`, ``, `name of the term`)
	updateCmd.Flags().StringP(`title`, `t`, ``, `title of the term`)
	updateCmd.Flags().String(`description`, ``, `description for this vocabulary`)

	deleteCmd := &cobra.Command{
		Use:   `delete [id]`,
		Args:  cobra.ExactArgs(1),
		Short: `Delete term`,
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.ParseUint(args[0], 10, 32)
			CheckErr(err)
			CheckErr(service(cmd).Term.Delete(cmd.Context(), uint(id)))
		},
	}

	listCmd := &cobra.Command{
		Use:   `list [vocabulary's id] [limit] [offset]`,
		Args:  cobra.MinimumNArgs(1),
		Short: `Show all terms`,
		Run: func(cmd *cobra.Command, args []string) {
			vocabularyID, err := strconv.ParseUint(args[0], 10, 32)
			CheckErr(err)
			limit, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				limit = 10
			}
			offset, err := strconv.ParseUint(args[2], 10, 32)
			if err != nil {
				offset = 0
			}
			terms, err := service(cmd).Term.Get(cmd.Context(), &model.TermFilter{
				VocabularyID: []uint64{uint64(vocabularyID)},
				Limit:        uint(limit),
				Offset:       uint64(offset),
			})
			CheckErr(err)

			table := tablewriter.NewWriter(cmd.OutOrStdout())
			table.SetHeader([]string{`ID`, `Name`, `Title`})

			for _, term := range terms {
				table.Append(func(term model.Term) []string {
					return []string{
						strconv.Itoa(int(term.ID)),
						term.Data.Name,
						term.Data.Title,
					}
				}(term))
			}
			table.Render()
		},
	}

	termCmd.AddCommand(createCmd, updateCmd, deleteCmd, listCmd)

	return termCmd
}
