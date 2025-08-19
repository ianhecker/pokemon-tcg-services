package networking

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type NetworkingClient interface {
	Get(url string) (body []byte, status int, err error)
}

const Timeout = 60 * time.Second

type Client struct {
	httpclient *http.Client
	logger     *zap.SugaredLogger
}

func NewClient(
	logger *zap.SugaredLogger,
	timeout time.Duration,
) NetworkingClient {
	return &Client{
		httpclient: &http.Client{
			Timeout: timeout,
		},
		logger: logger,
	}
}

func (client *Client) Get(url string) ([]byte, int, error) {
	log := client.logger
	log.Infow("client fetching url", "url", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {

		log.Errorw("client creating request", "url", url, "error", err)
		return nil, 0, fmt.Errorf("error creating request: %w", err)
	}

	start := time.Now()

	resp, err := client.httpclient.Do(req)
	if err != nil {

		log.Errorw("client fetching url", "url", url, "error", err)
		return nil, 0, fmt.Errorf("error doing request: %w", err)
	}
	defer resp.Body.Close()

	end := time.Now()
	elapsed := end.Sub(start)

	status := resp.StatusCode

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("error reading body: %w", err)
	}

	if status != http.StatusOK {
		log.Errorw("status not OK fetching url", "url", url, "time", elapsed, "status", status, "body", string(body))
		return body, status, fmt.Errorf("unexpected status code: %d", status)
	}

	log.Infow("fetched url", "url", url, "time", elapsed, "body", string(body))
	return body, status, nil
}
