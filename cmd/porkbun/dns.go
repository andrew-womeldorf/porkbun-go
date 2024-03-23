package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"

	"github.com/andrew-womeldorf/porkbun-go"
	"github.com/spf13/cobra"
)

var domain string

func initDnsCmd() {
	dnsCmd.AddCommand(dnsCreateCmd)

	dnsCmd.PersistentFlags().StringVarP(&domain, "domain", "d", "", "domain name")
	cobra.MarkFlagRequired(dnsCmd.PersistentFlags(), "domain")

	dnsCreateFlags := dnsCreateCmd.Flags()
	dnsCreateFlags.StringP("subdomain", "s", "", "subdomain for the record. Leave blank to create a record on the root domain")
	dnsCreateFlags.String("ttl", "600", "time to live for the record")
	dnsCreateFlags.String("priority", "", "priority of the record for those that support it")
	cobra.MarkFlagRequired(dnsCreateFlags, "type")
	cobra.MarkFlagRequired(dnsCreateFlags, "content")
}

var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "Manage DNS entries for a domain",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var dnsCreateCmd = &cobra.Command{
	Use:   "create TYPE CONTENT",
	Short: "Create a new DNS entry",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		priority, err := cmd.Flags().GetString("priority")
		if err != nil {
			log.Fatal(fmt.Errorf("err getting priority var, %w", err))
		}

		subdomain, err := cmd.Flags().GetString("subdomain")
		if err != nil {
			log.Fatal(fmt.Errorf("err getting subdomain var, %w", err))
		}

		ttl, err := cmd.Flags().GetString("ttl")
		if err != nil {
			log.Fatal(fmt.Errorf("err getting ttl var, %w", err))
		}

		req := &porkbun.CreateDnsRecordRequest{
			Name:     subdomain,
			Type:     args[0],
			Content:  args[1],
			TTL:      ttl,
			Priority: priority,
		}

		client, err := porkbun.NewClient()
		if err != nil {
			log.Fatal(fmt.Errorf("err creating porkbun client, %w", err))
		}

		slog.Debug("Sending request", "params", req, "domain", domain)

		res, err := client.CreateDnsRecord(ctx, domain, req)
		if err != nil {
			log.Fatal(fmt.Errorf("err creating dns record, %w", err))
		}

		resBytes, err := json.Marshal(res)
		if err != nil {
			log.Fatal(fmt.Errorf("error marshaling response to JSON, %w", err))
		}
		fmt.Println(string(resBytes))
	},
}
