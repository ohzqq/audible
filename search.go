package audible

import (
	"net/url"
	"slices"
)

var queryFields = []string{
	"keywords",
	"author",
	"narrator",
	"title",
}

func GetSearchParams(vals url.Values) url.Values {
	params := make(url.Values)
	for k, v := range vals {
		if slices.Contains(queryFields, k) {
			params[k] = v
		}
	}
	return params
}
