package cmd

import (
	"encoding/json"
	"log"
	"os"

	"github.com/ohzqq/audible"
	"github.com/ohzqq/audible/tui"
	"github.com/ohzqq/cdb"
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

		r, err := prod.Get()
		if err != nil {
			log.Fatal(err)
		}

		l := tui.New(r)
		books := l.Run()

		processProducts(books)
	},
}

func saveResults(products *audible.ProductsResponse) {
	d, err := json.Marshal(products)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("search-results.json", d, 0644)
	if err != nil {
		log.Fatal(err)
	}

}

func pickResult(products []audible.Product) {
	var items []string
	var books []cdb.Book
	for _, p := range products {
		b := p.ToBook()
		books = append(books, b)
		items = append(items, b.Title)
	}

}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.PersistentFlags().StringVarP(&flagAuthor, "author", "a", "", "search authors")
	searchCmd.PersistentFlags().StringVarP(&flagTitle, "title", "t", "", "search title")
	searchCmd.PersistentFlags().StringVarP(&flagNarrator, "narrator", "n", "", "search narrators")
}
