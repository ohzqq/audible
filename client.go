package audible

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/spf13/cast"
)

const (
	apiSub                = `api`
	webSub                = `www`
	domain                = `audible`
	tld                   = `com`
	apiHost               = `api.audible.`
	productResponseGroups = `media,product_desc,contributors,series,product_extended_attrs,product_attrs`
	contentResponseGroups = `chapter_info,content_reference`
	products              = `/1.0/catalog/products`
	content               = `/1.0/content`
	asinRegex             = `/(?P<end>[^/]+)/?(?P<slug>[^/]+)?\/(?P<asin>\w{10})$`
)

var client = &http.Client{}

var responseGroups = map[string]string{
	content:  contentResponseGroups,
	products: productResponseGroups,
}

type host []string

type Request struct {
	endpoint string
	Asin     string
	host     host
	params   url.Values
	search   url.Values
	*url.URL
	*Query
}

type Response struct {
	Products     []Product `json:"products"`
	TotalResults int       `json:"total_results"`
	Product      Product   `json:"product"`
}

func NewRequest() *Request {
	req := &Request{
		Query:  New(),
		search: make(url.Values),
	}
	req.SubDomain(apiSub)
	req.Domain(domain)
	req.TLD(tld)
	return req
}

func get(u string) ([]byte, error) {
	hr, err := client.Get(u)
	if err != nil {
		return []byte{}, err
	}

	body, err := io.ReadAll(hr.Body)
	if err != nil {
		return []byte{}, err
	}
	defer hr.Body.Close()

	return body, nil
}

func (req *Request) String() string {
	return req.Query.String()
}

func (req *Request) ParseURL(uri string) (*Request, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return req, err
	}

	asin, err := AsinFromPath(u.Path)
	if err != nil {
		return req, err
	}
	req.Asin = asin
	req.AppendPath(asin)

	req.TLD(MarketFromHost(u.Host))

	return req, nil
}

func (r *Request) Title(kw ...string) *Request {
	if len(kw) > 0 {
		r.AddParam("title", strings.Join(kw, " "))
	}
	return r
}

func (r *Request) Author(kw ...string) *Request {
	if len(kw) > 0 {
		r.AddParam("author", strings.Join(kw, " "))
	}
	return r
}

func (r *Request) Narrator(kw ...string) *Request {
	if len(kw) > 0 {
		r.AddParam("narrator", strings.Join(kw, " "))
	}
	return r
}

func (r *Request) Keywords(kw ...string) *Request {
	if len(kw) > 0 {
		r.AddParam("keywords", strings.Join(kw, " "))
	}
	return r
}

func (r *Request) EncodeSearch() string {
	return r.Query.Query.Encode()
}

func MarketFromHost(host string) string {
	return strings.TrimPrefix(path.Ext(host), ".")
}

func AsinFromPath(path string) (string, error) {
	var aR = regexp.MustCompile(asinRegex)
	parsed := aR.FindStringSubmatch(path)
	if idx := aR.SubexpIndex("asin"); len(parsed) >= idx {
		return parsed[idx], nil
	}
	return "", errors.New("malformed url")
}

func (r *Request) NumResults(n int) *Request {
	r.SetParam("num_results", cast.ToString(n))
	return r
}

func newHost() host {
	return host([]string{apiSub, domain, tld})
}

func (h host) Sub(s string) {
	h[0] = s
}

func (h host) Domain(s string) {
	h[1] = s
}

func (h host) TLD(s string) {
	h[2] = s
}

func (h host) String() string {
	return strings.Join(h, ".")
}
