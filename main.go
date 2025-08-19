package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.uber.org/zap"

	"github.com/ianhecker/pokemon-tcg-services/networking"
	"github.com/ianhecker/pokemon-tcg-services/networking/pokemontcgio"
	v2 "github.com/ianhecker/pokemon-tcg-services/networking/pokemontcgio/v2"
	"github.com/ianhecker/pokemon-tcg-services/retry"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	httpClient := networking.NewClient(sugar)

	api := v2.CardByID("xy1-1")
	url := fmt.Sprintf("https://%s", api)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	perAttemptTimeout := 60 * time.Second
	client := pokemontcgio.NewClient(httpClient, perAttemptTimeout)

	result, retryFunc := client.MakeRetryFunc(url)

	retries := 10
	backoff := 2 * time.Second
	retryable := retry.MakeRetryable(retries, backoff, retryFunc)

	retryErr := retry.Do(ctx, retryable)
	checkErr(sugar, retryErr)

	var out bytes.Buffer
	err := json.Indent(&out, result.Body, "", " ")
	checkErr(sugar, err)

	fmt.Println(out.String())
}

func checkErr(logger *zap.SugaredLogger, err error) {
	if err != nil {
		logger.Errorw("error", "err", err)
		log.Fatal()
	}
}
