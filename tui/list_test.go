package tui

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/ohzqq/audible"
)

func TestList(t *testing.T) {
	d, err := os.ReadFile("../testdata/search-results.json")
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	var products audible.ProductsResponse
	err = json.Unmarshal(d, &products)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	l := New(&products)
	books := l.Run()
	fmt.Printf("%#v\n", books)
}
