package networking

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const Timeout = 10 * time.Second

type Client struct {
	httpclient *http.Client
}

func NewClient() *Client {
	client := &http.Client{
		Timeout: Timeout,
	}
	return NewClientFromRaw(client)
}

func NewClientFromRaw(client *http.Client) *Client {
	return &Client{
		httpclient: client,
	}
}

func (client *Client) Get(url *url.URL) ([]byte, error) {
	if url == nil {
		return nil, errors.New("url is nil")
	}

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := client.httpclient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error doing request: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}
	return body, nil
}
