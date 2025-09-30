package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/ianhecker/pokemon-tcg-services/internal/justtcg"
	"github.com/ianhecker/pokemon-tcg-services/internal/services/cardpricer"
)

var cardPricerCmd = &cobra.Command{
	Use:   "card-pricer",
	Short: "Fetch pricing for a pokemon card by card ID",
	Long:  `A service that fetches pokemon card prices via card ID`,
	Run: func(cmd *cobra.Command, args []string) {
		port := PortToString(Port)
		RunCardService(port)
	},
}

func RunCardService(port string) {
	err := Config.FlightCheck()
	if err != nil {
		log.Fatal(err)
	}

	base, _ := zap.NewProduction()
	defer base.Sync()
	logger := base.Sugar()

	ctx := context.Background()
	client := justtcg.NewClient(logger, Config.JustTCG.APIKey)

	svc := cardpricer.NewService(logger, client, port)
	stop := svc.Start(ctx)
	defer stop()

	select {
	case <-ctx.Done():
		logger.Errorw("context", "err", ctx.Err())
	case <-svc.Done():
		logger.Errorw("server", "err", svc.Err())
	}
}
