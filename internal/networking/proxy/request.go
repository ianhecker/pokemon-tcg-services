package proxy

import "net/http"

type Request struct {
	Req *http.Request
	Err error
}
