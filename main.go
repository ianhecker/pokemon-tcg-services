package main

import (
	"log"

	"go.uber.org/zap"

	"github.com/ianhecker/pokemon-tcg-services/networking"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	client := networking.NewClient(sugar, networking.Timeout)

	_, _, err := client.Get("https://api.pokemontcg.io/v2/cards/xy1-1")
	checkErr(sugar, err)

	// fmt.Println(string(body))
}

func checkErr(logger *zap.SugaredLogger, err error) {
	if err != nil {
		logger.Errorw("error", "err", err)
		log.Fatal()
	}
}
