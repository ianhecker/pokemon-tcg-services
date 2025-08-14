package networking

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const Timeout = 60 * time.Second

type Client struct {
	httpclient *http.Client
	logger     *zap.SugaredLogger
}

func NewClient(logger *zap.SugaredLogger) *Client {
	return &Client{
		httpclient: &http.Client{
			Timeout: Timeout,
		},
		logger: logger,
	}
}

func (client *Client) Get(url string) ([]byte, error) {
	client.logger.Infow("fetching url", "url", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	start := time.Now()

	resp, err := client.httpclient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error doing request: %w", err)
	}
	defer resp.Body.Close()

	end := time.Now()
	elapsed := end.Sub(start)

	client.logger.Infow("fetched url", "url", url, "time", elapsed)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}
	client.logger.Infow("payload for url", "url", url, "body", string(body))

	return body, nil
}
