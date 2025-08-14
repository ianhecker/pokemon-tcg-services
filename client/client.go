package client

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const Timeout = 10 * time.Second

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	httpclient HTTPClient
}

func NewClient() *Client {
	client := &http.Client{
		Timeout: Timeout,
	}
	return NewClientFromRaw(client)
}

func NewClientFromRaw(client HTTPClient) *Client {
	return &Client{
		httpclient: client,
	}
}

func (client *Client) NewRequest(method string, url string) (*Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	return &Request{req: req}, nil
}

func (client *Client) Do(request *Request) (*Response, error) {
	resp, err := client.httpclient.Do(request.req)
	if err != nil {
		return nil, fmt.Errorf("error doing request: %w", err)
	}

	return &Response{resp: resp}, nil
}

func (client *Client) Get(url *url.URL) ([]byte, error) {
	if url == nil {
		return nil, errors.New("url is nil")
	}

	request, err := client.NewRequest(http.MethodGet, url.String())
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {

	}

	body, err := response.ReadBody()
	if err != nil {
		return nil, err
	}
	return body, nil
}
