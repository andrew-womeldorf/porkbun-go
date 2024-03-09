package porkbun_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andrew-womeldorf/porkbun-go"
)

func TestClientPing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"status": "SUCCESS", "yourIp": "127.0.0.1"}`)
	}))
	defer server.Close()

	ctx := context.TODO()

	client, _ := porkbun.NewClient(
		porkbun.WithApiKey("apikey"),
		porkbun.WithSecretKey("secretkey"),
		porkbun.WithBaseUrl(server.URL),
	)

	res, err := client.Ping(ctx)
	if err != nil {
		t.Fatalf("got %s, want nil", err)
	}

	if res.Status != "SUCCESS" {
		t.Errorf("got %s, want %s", res.Status, "SUCCESS")
	}

	if res.YourIP != "127.0.0.1" {
		t.Errorf("got %s, want %s", res.YourIP, "127.0.0.1")
	}
}
