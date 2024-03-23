package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/andrew-womeldorf/porkbun-go"
	"github.com/spf13/cobra"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "porkbun",
	Short: "Interact with porkbun",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping with authentication",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		client, err := porkbun.NewClient()
		if err != nil {
			log.Fatal(err)
		}

		res, err := client.Ping(ctx)
		if err != nil {
			log.Fatal(err)
		}

		resBytes, err := json.Marshal(res)
		if err != nil {
			log.Fatal(fmt.Errorf("error marshaling response to JSON, %w", err))
		}
		fmt.Println(string(resBytes))
	},
}

func initLogger() {
	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
}

func init() {
	cobra.OnInitialize(initLogger)

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Output verbose logs")
	rootCmd.AddCommand(dnsCmd)
	rootCmd.AddCommand(pingCmd)

	initDnsCmd()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
