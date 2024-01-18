package audible

import (
	"net/url"
	"strings"
)

type Params struct {
	params url.Values
}

func NewSearch(kw ...string) *Params {
	s := &Params{params: make(url.Values)}
	if len(kw) > 0 {
		s.params.Set("keywords", strings.Join(kw, " "))
	}
	return s
}

func (s *Params) Title(kw ...string) *Params {
	if len(kw) > 0 {
		s.params.Set("title", strings.Join(kw, " "))
	}
	return s
}

func (s *Params) Author(kw ...string) *Params {
	if len(kw) > 0 {
		s.params.Set("author", strings.Join(kw, " "))
	}
	return s
}

func (s *Params) Narrator(kw ...string) *Params {
	if len(kw) > 0 {
		s.params.Set("narrator", strings.Join(kw, " "))
	}
	return s
}

func (s *Params) Encode() string {
	return s.params.Encode()
}

func (s *Params) String() string {
	return s.params.Encode()
}
