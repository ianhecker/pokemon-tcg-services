package proxy

import "net/http"

const GET Method = http.MethodGet

type Method string

func (m Method) String() string {
	return string(m)
}
