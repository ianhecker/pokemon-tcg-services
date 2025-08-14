package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Response struct {
	resp *http.Response
}

func NewResponse(resp *http.Response) *Response {
	return &Response{
		resp: resp,
	}
}

func (r Response) ReadBody() ([]byte, error) {
	if r.resp == nil {
		return nil, errors.New("response is nil")
	}

	defer r.resp.Body.Close()

	body, err := io.ReadAll(r.resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}
	return body, nil
}
