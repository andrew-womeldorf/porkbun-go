package porkbun_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andrew-womeldorf/porkbun-go"
)

func TestDns(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"status": "SUCCESS", "id": 1234}`)
		}))
		defer server.Close()

		ctx := context.TODO()

		client, _ := porkbun.NewClient(
			porkbun.WithApiKey("apikey"),
			porkbun.WithSecretKey("secretkey"),
			porkbun.WithBaseUrl(server.URL),
		)

		res, err := client.CreateDnsRecord(ctx, "example.com", &porkbun.Record{
			Name:    "www",
			Type:    "A",
			Content: "10.0.0.1",
		})
		if err != nil {
			t.Fatalf("got %s, want nil", err)
		}

		if res.Status != "SUCCESS" {
			t.Errorf("got %s, want %s", res.Status, "SUCCESS")
		}

		if res.Id != 1234 {
			t.Errorf("got %d, want %d", res.Id, 1234)
		}
	})
}
