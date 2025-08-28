package cmd

import (
	"log"

	"github.com/ianhecker/pokemon-tcg-services/internal/config"
	"github.com/spf13/cobra"
)

var Config config.Config

var rootCmd = &cobra.Command{
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		Config, err = config.Load()
		if err != nil {
			return err
		}
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(serviceCmd)
}
