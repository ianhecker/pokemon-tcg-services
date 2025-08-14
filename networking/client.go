package networking

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const Timeout = 10 * time.Second

type Client struct {
	httpclient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpclient: &http.Client{
			Timeout: Timeout,
		},
	}
}

func (client *Client) Get(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := client.httpclient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error doing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}
	return body, nil
}
