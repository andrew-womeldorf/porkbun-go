package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"strings"

	"github.com/andrew-womeldorf/porkbun-go"
	"github.com/spf13/cobra"
)

func initDnsCmd() {
	dnsCmd.AddCommand(dnsCreateCmd)
	dnsCmd.AddCommand(dnsListCmd)

	dnsCreateFlags := dnsCreateCmd.Flags()
	dnsCreateFlags.String("ttl", "600", "time to live for the record")
	dnsCreateFlags.String("priority", "", "priority of the record for those that support it")
}

var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "Manage DNS entries for a domain",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var dnsCreateCmd = &cobra.Command{
	Use:   "create DOMAIN TYPE CONTENT",
	Short: "Create a new DNS entry",
	Long: `Create a new DNS entry.

DOMAIN is the complete domain, such as 'foo.example.com', where 'foo' is the record entry on the 'example.com' domain.
TYPE is the type of record being created, such as A, AAAA, TXT, MX...
CONTENT is the answer for the record.`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		priority, err := cmd.Flags().GetString("priority")
		if err != nil {
			log.Fatal(fmt.Errorf("err getting priority var, %w", err))
		}

		ttl, err := cmd.Flags().GetString("ttl")
		if err != nil {
			log.Fatal(fmt.Errorf("err getting ttl var, %v", err))
		}

		sub, dom, err := ParseDomain(args[0])
		if err != nil {
			log.Fatal(fmt.Errorf("err parsing domain, %v", err))
		}

		req := &porkbun.CreateDnsRecordRequest{
			Name:     sub,
			Type:     args[1],
			Content:  args[2],
			TTL:      ttl,
			Priority: priority,
		}

		client, err := porkbun.NewClient()
		if err != nil {
			log.Fatal(fmt.Errorf("err creating porkbun client, %w", err))
		}

		slog.Debug("Sending create request", "params", req, "domain", dom)

		res, err := client.CreateDnsRecord(ctx, dom, req)
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

var dnsListCmd = &cobra.Command{
	Use:   "list DOMAIN",
	Short: "List entries for a domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		_, dom, err := ParseDomain(args[0])
		if err != nil {
			log.Fatal(fmt.Errorf("err parsing domain, %v", err))
		}

		client, err := porkbun.NewClient()
		if err != nil {
			log.Fatal(fmt.Errorf("err creating porkbun client, %w", err))
		}

		slog.Debug("Sending list request", "domain", dom)

		res, err := client.ListDnsRecords(ctx, dom)
		if err != nil {
			log.Fatal(fmt.Errorf("err listing dns records, %w", err))
		}

		resBytes, err := json.Marshal(res)
		if err != nil {
			log.Fatal(fmt.Errorf("error marshaling response to JSON, %w", err))
		}
		fmt.Println(string(resBytes))
	},
}

// ParseDomain takes a full domain as an input, and return the subdomain, the
// domain, and an error.
//
// "example.com" -> "", "example.com", nil
// "foo.example.com" -> "foo", "example.com", nil
// "*.example.com" -> "*", "example.com", nil
// "foo.bar.example.com" -> "foo.bar", "example.com", nil
func ParseDomain(domain string) (string, string, error) {
	// Split on "."
	parts := strings.Split(domain, ".")

	// Handle error cases
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid domain %s", domain)
	}

	// The domain is the last two parts
	dom := strings.Join(parts[len(parts)-2:], ".")

	// The subdomain is everything before the last two parts
	var sub string
	if len(parts) > 2 {
		sub = strings.Join(parts[:len(parts)-2], ".")
	}

	return sub, dom, nil
}
