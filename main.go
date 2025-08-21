package main

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/ianhecker/pokemon-tcg-services/networking"
	v2 "github.com/ianhecker/pokemon-tcg-services/pokemontcgio/v2"
	"github.com/ianhecker/pokemon-tcg-services/retry"
)

func main() {
	base, _ := zap.NewProduction()
	defer base.Sync()
	logger := base.Sugar()

	httpClient := networking.NewClient(logger)

	api := v2.CardByID("xy1-1")
	url := fmt.Sprintf("https://%s", api)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	perAttemptTimeout := 60 * time.Second
	client := v2.NewClient(logger, httpClient, perAttemptTimeout)

	_, retryFunc := client.MakeRetryFunc(url)

	retries := 10
	backoff := 2 * time.Second
	retryable := retry.MakeRetryable(retries, backoff, retryFunc)

	err := retry.RunRetryable(ctx, retryable)
	if err != nil {
		logger.Errorw("retryable error", "url", url, "err", err)
		return
	}
}
