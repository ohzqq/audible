package cmd

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/danielgtaylor/casing"
	"github.com/ohzqq/audbk"
	"github.com/ohzqq/audible"
	"github.com/ohzqq/avtools/cue"
	"github.com/spf13/cobra"
)

var (
	flagURL      string
	noMeta       bool
	noCover      bool
	withChapters bool
	flagExt      string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "urmeta",
	Short: "scrape book metadata",
	Long:  `scrape book/work metadata from audible and ao3`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&flagURL, "url", "u", "", "scrape url")
	rootCmd.PersistentFlags().StringVarP(&flagExt, "ext", "e", ".yaml", "extension to save")
	rootCmd.PersistentFlags().BoolVar(&noMeta, "no-meta", false, "don't write metadata to disk")
	rootCmd.PersistentFlags().BoolVar(&noCover, "no-cover", false, "don't write cover to disk")
	rootCmd.PersistentFlags().BoolVarP(&withChapters, "with-chapters", "c", false, "don't write cover to disk")
}

func processProducts(prods []audible.Product) {
	for _, prod := range prods {
		book := prod.ToBook()
		println(book.Title)
		name := casing.Snake(book.Title)
		if !noMeta {
			book.Save(name+flagExt, true)
			ff := audbk.NewFFMeta()
			audbk.BookToFFMeta(ff, book.StringMap())
			ffm, err := os.Create(name + ".ini")
			if err != nil {
				log.Fatal(err)
			}
			defer ffm.Close()
			ff.WriteTo(ffm)
		}
		if !noCover {
			getCover(name+".jpg", book.Cover)
		}
		if withChapters {
			getChaps(name+".cue", prod.Asin)
		}
	}
}

func getChaps(name, asin string) {
	res, err := audible.Content().Asin(asin)
	if err != nil {
		log.Fatal(err)
	}
	d := cue.Dump("x.m4b", res.Chapters())
	err = os.WriteFile(name+".cue", d, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func getCover(n, u string) {
	res, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	file, err := os.Create(n)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		log.Fatal(err)
	}
}
