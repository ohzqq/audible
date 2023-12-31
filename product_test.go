package audible

import (
	"fmt"
	"testing"
)

type searchTest struct {
	kw   []string
	want string
}

var searchTests = []searchTest{
	searchTest{
		kw:   []string{},
		want: "",
	},
	searchTest{
		kw:   []string{"red fish"},
		want: "keywords=red+fish",
	},
	searchTest{
		kw:   []string{"red", "fish"},
		want: "keywords=red+fish",
	},
}

func TestNewSearch(t *testing.T) {
	for _, test := range searchTests {
		search := NewSearch(test.kw...)
		if u := search.Encode(); u != test.want {
			t.Errorf("output %#v != expected %v", u, test.want)
		}
	}
}

var paramTests = map[string]searchTest{
	"author": searchTest{
		kw:   []string{"amy", "lane"},
		want: "author=amy+lane",
	},
	"narrator": searchTest{
		kw:   []string{"greg", "tremblay"},
		want: "narrator=greg+tremblay",
	},
	"title": searchTest{
		kw:   []string{"red", "fish"},
		want: "title=red+fish",
	},
}

var combinedQuery = "author=amy+lane&narrator=greg+tremblay&title=red+fish"

func TestSearchParams(t *testing.T) {
	for name, test := range paramTests {
		search := NewSearch()
		switch name {
		case "title":
			search.Title(test.kw...)
		case "author":
			search.Author(test.kw...)
		case "narrator":
			search.Narrator(test.kw...)
		}
		if u := search.Encode(); u != test.want {
			t.Errorf("output %#v != expected %v", u, test.want)
		}
	}
}

func TestCombinedParams(t *testing.T) {
	search := NewSearch()
	for name, test := range paramTests {
		if name == "author" {
			search.Author(test.kw...)
		}
		if name == "narrator" {
			search.Narrator(test.kw...)
		}
		if name == "title" {
			search.Title(test.kw...)
		}
	}
	if u := search.Encode(); u != combinedQuery {
		t.Errorf("output %#v != expected %v", u, combinedQuery)
	}
}

var testURLs = []string{
	"https://www.audible.com/pd/Red-Fish-Dead-Fish-Audiobook/B07B4LFT72",
	"https://www.audible.com/pd/Red-Fish-Dead-Fish-Audiobook/B07B4LT2",
	"https://www.audible.com/series/Fish-out-of-Water-Audiobook/B07B5HD42Y",
}

func TestGetFromURL(t *testing.T) {
	for _, test := range testURLs {
		r, err := Products().URL(test)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		if len(r.Products) < 1 {
			fmt.Printf("has %d results, expected at least one\n", len(r.Products))
		}
	}
}

func TestSearchResults(t *testing.T) {
	for name, test := range paramTests {
		search := NewSearch()
		switch name {
		case "title":
			search.Title(test.kw...)
		case "author":
			search.Author(test.kw...)
		case "narrator":
			search.Narrator(test.kw...)
		}
		p := Products()
		p.NumResults(1)
		_, err := p.Search(search)
		if err != nil {
			t.Error(err)
		}
		//println(search.Encode())
		//fmt.Printf("%#v\n", r)
	}

	for _, test := range searchTests {
		search := NewSearch(test.kw...)
		p := Products()
		p.NumResults(1)
		_, err := p.Search(search)
		if err != nil {
			t.Error(err)
		}
		//println(search.Encode())
		//fmt.Printf("%#v\n", r)
	}
}
