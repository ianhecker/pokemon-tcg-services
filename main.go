package main

import (
	"context"

	v2 "github.com/ianhecker/pokemon-tcg-services/internal/pokemontcgio/v2"
	"github.com/ianhecker/pokemon-tcg-services/internal/services/card"
	"go.uber.org/zap"
)

func main() {
	base, _ := zap.NewProduction()
	defer base.Sync()
	logger := base.Sugar()

	ctx := context.Background()
	client := v2.NewClient(logger)

	svc := card.NewService(logger, client, ":8080")
	stop := svc.Start(ctx)
	defer stop()

	select {
	case <-ctx.Done():
		logger.Errorw("context", "err", ctx.Err())
	case <-svc.Done():
		logger.Errorw("server", "err", svc.Err())
	}
}
