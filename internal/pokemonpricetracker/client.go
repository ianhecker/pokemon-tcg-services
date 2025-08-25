package pokemonpricetracker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ianhecker/pokemon-tcg-services/internal/networking"
	"github.com/ianhecker/pokemon-tcg-services/internal/pokemontcg"
	"github.com/ianhecker/pokemon-tcg-services/internal/retry"
	"go.uber.org/zap"
)

const Retries int = 5
const BackoffInSeconds = 1 * time.Second

type APIClientInterface interface {
	GetPricing(ctx context.Context, card *pokemontcg.Card) error
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

func NewClient(logger *zap.SugaredLogger) *Client {
	httpClient := networking.NewClient(logger)

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

func (c *Client) GetPricing(ctx context.Context, card *pokemontcg.Card) error {
	if card == nil {
		return errors.New("card is nil")
	}
	url := MakePricesAPI(*card)

	result, retryFunc := c.MakeRetryFunc(url)
	retryable := retry.MakeRetryable(Retries, BackoffInSeconds, retryFunc)

	err := retry.RunRetryable(ctx, retryable)
	if err != nil {
		return err
	}

	err = json.Unmarshal(result.Body, card)
	if err != nil {
		return fmt.Errorf("error unmarshaling into card: %w", err)
	}
	return nil
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
