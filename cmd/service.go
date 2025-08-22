package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Port uint

func PortToString(n uint) string {
	return fmt.Sprintf(":%d", n)
}

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Run a microservice",
}

func init() {
	serviceCmd.PersistentFlags().UintVar(&Port, "port", 0, "port to bind")
	serviceCmd.MarkPersistentFlagRequired("port")

	serviceCmd.AddCommand(cardByIDServiceCmd)
}
