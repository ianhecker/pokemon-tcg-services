package client

import (
	"net/http"
	"net/url"
)

type Request struct {
	req *http.Request
}

func (r Request) Method() string {
	return r.req.Method
}

func (r Request) URL() *url.URL {
	return r.req.URL
}
