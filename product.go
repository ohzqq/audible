package audible

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/ohzqq/cdb"
	"github.com/spf13/cast"
)

type ProductsRequest struct {
	*Request
}

type ProductsResponse struct {
	Products     []Product `json:"products"`
	TotalResults int       `json:"total_results"`
	Product      Product   `json:"product"`
}

type Product struct {
	Authors          []map[string]string `json:"authors"`
	ProductImages    map[string]string   `json:"product_images"`
	PublisherSummary string              `json:"publisher_summary"`
	RuntimeMin       int                 `json:"runtime_length_min"`
	Asin             string              `json:"asin"`
	Languages        string              `json:"language"`
	Narrators        []map[string]string `json:"narrators"`
	Publisher        string              `json:"publisher_name"`
	ReleaseDate      string              `json:"release_date"`
	Series           []map[string]string `json:"series"`
	Title            string              `json:"title"`
	ChapterInfo
}

func Products() *ProductsRequest {
	req := &ProductsRequest{Request: newRequest()}
	req.SetParam("response_groups", responseGroups[products])
	req.AppendPath(products)
	req.NumResults(50)
	return req
}

func (p *ProductsRequest) Search(s url.Values) (*ProductsResponse, error) {
	for k, v := range s {
		for _, a := range v {
			p.AddParam(k, a)
		}
	}
	return p.Get()
}

func (p *ProductsRequest) URL(u string) (*ProductsResponse, error) {
	_, err := p.ParseURL(u)
	if err != nil {
		return &ProductsResponse{}, err
	}

	r, err := p.Get()
	if err != nil {
		return r, err
	}

	r.Products = []Product{r.Product}
	r.Product = Product{}

	return r, nil
}

func (p *ProductsRequest) Get() (*ProductsResponse, error) {
	res := &ProductsResponse{}

	d, err := get(p.String())
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(d, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}
func (p Product) ToBook() cdb.Book {
	book := cdb.Book{
		EditableFields: cdb.EditableFields{
			Title:       p.Title,
			Publisher:   p.Publisher,
			Languages:   []string{p.Languages},
			Comments:    p.PublisherSummary,
			Identifiers: []string{"asin:" + p.Asin},
		},
	}

	for _, a := range p.Authors {
		book.Authors = append(book.Authors, a["name"])
	}

	for _, n := range p.Narrators {
		book.Narrators = append(book.Narrators, n["name"])
	}

	if len(p.Series) > 0 {
		book.Series = p.Series[0]["title"]
		book.SeriesIndex = cast.ToFloat64(p.Series[0]["sequence"])
	}

	if href, ok := p.ProductImages["500"]; ok {
		book.Cover = href
	}

	h := p.RuntimeMin / 60
	m := p.RuntimeMin % 60
	book.Duration = fmt.Sprintf("%02d:%02d:%02d", h, m, 0)

	t, err := time.Parse(time.DateOnly, p.ReleaseDate)
	if err != nil {
		t = time.Now()
	}
	book.Pubdate = t

	return book
}

func (p Product) FilterValue() string {
	var auths []string
	for _, a := range p.Authors {
		auths = append(auths, a["name"])
	}
	auth := strings.Join(auths, " & ")
	return fmt.Sprintf("%s %s", p.Title, auth)
}

//func (p Product)

func parseAudibleSeries(series map[string]string) (string, float64) {
	var title string
	var pos float64
	for k, v := range series {
		if k == "sequence" {
			pos = cast.ToFloat64(v)
		}
		if k == "title" {
			title = v
		}
	}
	return title, pos
}
