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

type HttpClientInterface interface {
	Get(ctx context.Context, url string) (body []byte, status int, err error)
}

type Client struct {
	httpclient *http.Client
	logger     *zap.SugaredLogger
}

func NewClient(
	logger *zap.SugaredLogger,
) HttpClientInterface {
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
	log.Infow("requesting", "url", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {

		log.Errorw("creating request", "url", url, "error", err)
		return nil, 0, fmt.Errorf("error creating request: %w", err)
	}

	duration, body, status, err := client.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Infow("request canceled",
				"url", url,
				"time", duration.String(),
				"error", err,
			)
		} else if errors.Is(err, context.DeadlineExceeded) {
			log.Infow("request deadline exceeded",
				"url", url,
				"time", duration.String(),
				"error", err,
			)
		} else {
			log.Errorw("request error",
				"url", url,
				"time", duration.String(),
				"error", err,
			)
		}
		return nil, status, fmt.Errorf("error doing request: %w", err)
	}

	if status != http.StatusOK {
		log.Warnw("unexpected status code",
			"url", url,
			"time", duration.String(),
			"status", status,
			"body", string(body),
		)
		return body, status, fmt.Errorf("unexpected status code: %d", status)
	}

	log.Infow("got response",
		"url", url,
		"time", duration.String(),
		"body", string(body))
	return body, status, nil
}

func (client *Client) Do(request *http.Request) (time.Duration, []byte, int, error) {
	start := time.Now()

	resp, err := client.httpclient.Do(request)
	if err != nil {
		elapsed := time.Since(start)
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
