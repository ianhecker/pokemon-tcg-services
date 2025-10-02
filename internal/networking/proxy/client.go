package proxy

import "net/http"

type ClientInterface interface {
	Do(*http.Request) (*http.Response, error)
}
