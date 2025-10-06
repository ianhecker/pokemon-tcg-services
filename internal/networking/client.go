package networking

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/ianhecker/pokemon-tcg-services/internal/config"
	"github.com/ianhecker/pokemon-tcg-services/internal/networking/proxy"
)

type ClientInterface interface {
	Get(ctx context.Context, url string) (body []byte, status int, err error)
}

type Client struct {
	log   Logger
	proxy proxy.ProxyInterface
	token config.Token
}

func NewClient(
	logger *zap.SugaredLogger,
	token config.Token,
) ClientInterface {
	return NewClientFromRaw(
		MakeLogger(logger),
		proxy.NewProxy(),
		token,
	)
}

func NewClientFromRaw(
	logger Logger,
	proxy proxy.ProxyInterface,
	token config.Token,
) ClientInterface {
	return &Client{
		log:   logger,
		proxy: proxy,
		token: token,
	}
}

func (client *Client) Get(ctx context.Context, url string) ([]byte, int, error) {
	log := client.log
	log.Set(url)
	log.Requesting()

	req := client.proxy.NewRequest(ctx, proxy.GET, url)
	if req.Err != nil {
		log.RequestError(req.Err)
		return nil, 0, fmt.Errorf("client: %w", req.Err)
	}
	client.proxy.SetAuthorization(req, "X-API-Key", client.token.Reveal())

	resp := client.proxy.Do(req)
	elapsed := resp.Timer.Elapsed()
	body, status, err := resp.Body, resp.Status, resp.Err

	if resp.Err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			log.ContextIssue(elapsed, err)
		} else {
			log.ResponseError(elapsed, err)
		}
		return nil, status, fmt.Errorf("client: %w", err)
	}

	if status != http.StatusOK {
		log.UnexpectedStatus(elapsed, status, err)
		return body, status, fmt.Errorf("client: unexpected status code: %d", status)
	}
	log.Success(elapsed)
	return body, status, nil
}
