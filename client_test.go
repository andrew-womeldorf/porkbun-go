package porkbun_test

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/andrew-womeldorf/porkbun-go"
)

func TestClientSetBaseUrl(t *testing.T) {
	_, err := porkbun.NewClient(porkbun.WithBaseUrl("http://localhost:3000"))
	if err != nil {
		t.Fatal("did not create porkbun client")
	}
}

func TestClientSetHttpClient(t *testing.T) {
	_, err := porkbun.NewClient(porkbun.WithHttpClient(&http.Client{
		Timeout: 30 * time.Second,
	}))
	if err != nil {
		t.Fatal("did not create porkbun client")
	}
}

// Use the Ping method to verify that access keys are set properly, and errors
// correctly.
func TestAccessKeyHandlers(t *testing.T) {
	apiKey := "apikey"
	secretKey := "secretapikey"

	testCases := []struct {
		msg       string
		apiKey    string
		secretKey string
		missing   string
	}{
		{
			msg:       "missing api key",
			apiKey:    "",
			secretKey: "",
			missing:   porkbun.PORKBUN_API_KEY,
		},
		{
			msg:       "missing secret key",
			apiKey:    apiKey,
			secretKey: "",
			missing:   porkbun.PORKBUN_SECRET_KEY,
		},
		{
			msg:       "access keys present",
			apiKey:    apiKey,
			secretKey: secretKey,
			missing:   "",
		},
	}
	for _, tc := range testCases {
		t.Run("missing api key", func(t *testing.T) {
			ctx := context.TODO()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Fatal(err)
				}

				want := fmt.Sprintf(`{"apikey":"%s","secretapikey":"%s"}`, apiKey, secretKey)

				if string(body) != want {
					t.Fatalf("got %s, want %s", string(body), want)
				}
			}))

			client, _ := porkbun.NewClient(
				porkbun.WithApiKey(tc.apiKey),
				porkbun.WithSecretKey(tc.secretKey),
				porkbun.WithBaseUrl(server.URL),
			)

			_, err := client.Ping(ctx)

			if tc.missing != "" {
				if err == nil {
					t.Errorf("expected error")
				}

				var got porkbun.MissingAccessKeyError
				isMissingAccessKeyError := errors.As(err, &got)
				want := porkbun.MissingAccessKeyError{Key: tc.missing}

				if !isMissingAccessKeyError {
					t.Fatalf("was not a MissingAccessKeyError, got %T", err)
				}

				if got != want {
					t.Errorf("got %v, want %v", got, want)
				}
			}
		})
	}
}
