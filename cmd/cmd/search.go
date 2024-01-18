package cmd

import (
	"log"

	"github.com/ohzqq/audible"
	"github.com/ohzqq/srch"
	"github.com/ohzqq/srch/ui"
	"github.com/spf13/cobra"
)

var (
	flagAuthor   string
	flagTitle    string
	flagNarrator string
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:     "search keywords",
	Short:   "search audible api",
	Aliases: []string{"s"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prod := audible.Products()
		prod.Keywords(args...)
		if flagAuthor != "" {
			prod.Author(flagAuthor)
		}
		if flagNarrator != "" {
			prod.Narrator(flagNarrator)
		}
		if flagTitle != "" {
			prod.Title(flagTitle)
		}

		idx := srch.New("field=title").Index(prod.Search())

		app := ui.New(idx)
		sel, err := app.Run()
		if err != nil {
			log.Fatal(err)
		}

		if sel == nil {
			println("no results.")
			return
		}

		processResults(sel.Data...)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.PersistentFlags().StringVarP(&flagAuthor, "author", "a", "", "search authors")
	searchCmd.PersistentFlags().StringVarP(&flagTitle, "title", "t", "", "search title")
	searchCmd.PersistentFlags().StringVarP(&flagNarrator, "narrator", "n", "", "search narrators")
}
