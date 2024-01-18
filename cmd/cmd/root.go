package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/danielgtaylor/casing"
	"github.com/ohzqq/audbk"
	"github.com/ohzqq/audible"
	"github.com/ohzqq/avtools/cue"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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
	Use:   "audible",
	Short: "scrape audible metadata",
	Long:  `scrape book/work metadata from audible`,
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

func processResults(res ...map[string]any) {
	for _, r := range res {
		if !noMeta {
			println(r["title"].(string))
			name := casing.Snake(r["title"].(string))

			err := writeMetaFile(r, name)
			if err != nil {
				log.Fatal(err)
			}

			err = writeFFMeta(r, name)
			if err != nil {
				log.Fatal(err)
			}

			if !noCover {
				if c, ok := r["cover"]; ok {
					getCover(name+".jpg", c.(string))
				}
			}

			if withChapters {
				if ident, ok := r["identifiers"]; ok {
					var asin string
					for _, id := range ident.([]string) {
						if strings.HasPrefix(id, "asin:") {
							asin = strings.TrimPrefix(id, "asin:")
						}
					}
					getChaps(name+".cue", asin)
				}
			}
		}
	}
}

func writeFFMeta(r map[string]any, name string) error {
	ff := audbk.NewFFMeta()
	audbk.BookToFFMeta(ff, r)

	ffm, err := os.Create(name + ".ini")
	if err != nil {
		return fmt.Errorf("write init error: %w\n", err)
	}
	defer ffm.Close()
	ff.WriteTo(ffm)

	return nil
}

func writeMetaFile(r map[string]any, name string) error {
	var err error

	mf, err := os.Create(name + flagExt)
	defer mf.Close()

	switch flagExt {
	case ".yaml":
		err = yaml.NewEncoder(mf).Encode(r)
	case ".json":
		err = json.NewEncoder(mf).Encode(r)
	case ".toml":
		err = toml.NewEncoder(mf).Encode(r)
	}

	if err != nil {
		return fmt.Errorf("write meta file err %w\n", err)
	}
	return nil
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
