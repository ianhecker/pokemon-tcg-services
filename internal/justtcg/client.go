package justtcg

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/ianhecker/pokemon-tcg-services/internal/config"
	"github.com/ianhecker/pokemon-tcg-services/internal/justtcg/v1/cards"
	"github.com/ianhecker/pokemon-tcg-services/internal/networking"
	"github.com/ianhecker/pokemon-tcg-services/internal/retry"
)

const Retries int = 5
const BackoffInSeconds = 1 * time.Second

type APIClientInterface interface {
	GetPricing(ctx context.Context, ID cards.TCGPlayerID) (cards.Card, error)
}

type Result struct {
	Body   []byte
	Status int
	Err    error
}

type Client struct {
	log     *zap.SugaredLogger
	client  networking.HttpClientInterface
	timeout time.Duration
}

func NewClient(
	logger *zap.SugaredLogger,
	token config.Token,
) *Client {
	httpClient := networking.NewClient(logger, token)
	return &Client{
		log:    logger,
		client: httpClient,
	}
}

func NewClientFromRaw(
	logger *zap.SugaredLogger,
	client networking.HttpClientInterface,
) *Client {
	return &Client{
		log:    logger,
		client: client,
	}
}

func (c *Client) GetPricing(ctx context.Context, ID cards.TCGPlayerID) (cards.Card, error) {
	url := cards.GetCardByID(ID)
	result, retryFunc := c.MakeRetryFunc(url)
	retryable := retry.MakeRetryable(Retries, BackoffInSeconds, retryFunc)

	err := retry.RunRetryable(ctx, retryable)
	if err != nil {
		return cards.Card{}, err
	}

	response, err := cards.Decode(result.Body, cards.UseNumber())
	if err != nil {
		return cards.Card{}, err
	}
	card, err := cards.Map(response)
	if err != nil {
		return cards.Card{}, err
	}
	return card, nil
}

func (c *Client) MakeRetryFunc(url string) (*Result, retry.RetryFunc) {

	result := &Result{}
	retryFunc := func(ctx context.Context) (retry.RetryState, error) {

		body, status, err := c.client.Get(ctx, url)
		result.Body = body
		result.Status = status
		result.Err = err

		if err != nil {
			return retry.Fail, err
		}

		if status != 0 {
			state := RetryForStatus(status)
			if state == retry.Complete {
				return state, nil
			} else {
				return state, err
			}
		}
		return retry.Fail, nil
	}
	return result, retryFunc
}
