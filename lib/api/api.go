package api

import (
	"net/http"
	"net/url"
)

var API_ROOT = url.URL{
	Scheme: "https",
	Host:   "api.applyfinch.com",
	Path:   "/api/v1",
}

func Get(path string, query *url.Values) (*http.Response, error) {
	url := API_ROOT.Clone()
	url.Path += path
	url.RawQuery = query.Encode()
	return http.Get(url.String())
}
