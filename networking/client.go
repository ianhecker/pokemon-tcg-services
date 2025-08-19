package networking

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type ClientInterface interface {
	Get(ctx context.Context, url string) (body []byte, status int, err error)
}

type Client struct {
	httpclient *http.Client
	logger     *zap.SugaredLogger
}

func NewClient(
	logger *zap.SugaredLogger,
) ClientInterface {
	return &Client{
		httpclient: &http.Client{
			Timeout:   0,
			Transport: NewTransport(),
		},
		logger: logger,
	}
}

func (client *Client) Get(ctx context.Context, url string) ([]byte, int, error) {
	log := client.logger
	log.Infow("client fetching url", "url", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {

		log.Errorw("client creating request", "url", url, "error", err)
		return nil, 0, fmt.Errorf("error creating request: %w", err)
	}

	duration, body, status, err := client.Do(req)
	if err != nil {
		return nil, status, fmt.Errorf("error doing request: %w", err)
	}

	if status != http.StatusOK {
		log.Warnw("unexpected status code",
			"url", url,
			"time", duration,
			"status", status,
			"body", string(body),
		)
		return body, status, fmt.Errorf("unexpected status code: %d", status)
	}

	log.Infow("fetched url", "url", url, "time", duration, "body", string(body))
	return body, status, nil
}

func (client *Client) Do(request *http.Request) (time.Duration, []byte, int, error) {
	log := client.logger
	start := time.Now()

	resp, err := client.httpclient.Do(request)
	if err != nil {
		elapsed := time.Since(start)

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			log.Infow("request canceled via context", "url", request.URL.String(), "time", elapsed, "error", err)

		} else {
			log.Errorw("client request error", "url", request.URL.String(), "time", elapsed, "error", err)
		}
		return elapsed, nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	elapsed := time.Since(start)

	if err != nil {
		return elapsed, nil, resp.StatusCode, fmt.Errorf("error reading body: %w", err)
	}
	return elapsed, body, resp.StatusCode, nil
}
