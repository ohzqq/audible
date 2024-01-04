package cmd

import (
	"log"

	"github.com/ohzqq/audible"
	"github.com/spf13/cobra"
)

// urlCmd represents the url command
var urlCmd = &cobra.Command{
	Use:     "url url",
	Aliases: []string{"u"},
	Short:   "scrape audible url",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		req := audible.Products()
		r, err := req.URL(args[0])
		if err != nil {
			log.Fatal(err)
		}
		processProducts(r.Products)
	},
}

func init() {
	rootCmd.AddCommand(urlCmd)
}
