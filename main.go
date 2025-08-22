package main

import (
	"context"
	"fmt"
	"time"

	v2 "github.com/ianhecker/pokemon-tcg-services/internal/pokemontcgio/v2"
	"github.com/ianhecker/pokemon-tcg-services/internal/retry"
	"github.com/ianhecker/pokemon-tcg-services/internal/services/card"
	"go.uber.org/zap"
)

func main() {
	base, _ := zap.NewProduction()
	defer base.Sync()
	logger := base.Sugar()

	api := v2.CardByID("xy1-1")
	url := fmt.Sprintf("https://%s", api)

	client := v2.NewClient(logger)
	_, retryable := client.MakeRetryable(url)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	err := retry.RunRetryable(ctx, retryable)
	if err != nil {
		logger.Errorw("retryable error", "url", url, "err", err)
		return
	}

	svc := card.NewService(logger, nil, ":8080")
	stop := svc.Start(ctx)
	defer stop()

	select {
	case <-ctx.Done():
		logger.Errorw("context", "err", ctx.Err())
	case <-svc.Done():
		logger.Errorw("server", "err", svc.Err())
	}
}
