package porkbun_test

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
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

func echoServer(t testing.TB) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		w.Write(body)
	}))
}

func TestClientDo(t *testing.T) {
	testCases := []struct {
		msg      string
		body     io.Reader
		expected string
	}{
		{
			msg:      "nil body does not error",
			body:     nil,
			expected: `{"apiKey":"apikey","secretKey":"secretkey"}`,
		},
		{
			msg:      "nil body does not error",
			body:     strings.NewReader(`{"foo": "bar"}`),
			expected: `{"apiKey":"apikey","foo":"bar","secretKey":"secretkey"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			server := echoServer(t)
			defer server.Close()

			client, err := porkbun.NewClient(
				porkbun.WithApiKey("apikey"),
				porkbun.WithSecretKey("secretkey"),
			)
			if err != nil {
				t.Fatal(fmt.Errorf("did not create porkbun client, %w", err))
			}

			req, err := http.NewRequest(http.MethodGet, server.URL, tc.body)
			if err != nil {
				t.Fatal(fmt.Errorf("did not create new request, %w", err))
			}

			res, err := client.Do(req)
			if err != nil {
				t.Fatal(fmt.Errorf("client.Do should not error, %w", err))
			}
			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatal(fmt.Errorf("coult not read response body, %w", err))
			}

			if string(body) != tc.expected {
				t.Errorf("got %s, want %s", string(body), tc.expected)
			}
		})
	}

	testCasesAccessKeys := []struct {
		apiKey    string
		secretKey string
		missing   string
	}{
		{
			apiKey:    "",
			secretKey: "",
			missing:   porkbun.PORKBUN_API_KEY,
		},
		{
			apiKey:    "foo",
			secretKey: "",
			missing:   porkbun.PORKBUN_SECRET_KEY,
		},
	}
	for _, tc := range testCasesAccessKeys {
		t.Run("missing api key", func(t *testing.T) {
			server := echoServer(t)
			defer server.Close()

			client, _ := porkbun.NewClient(porkbun.WithApiKey(tc.apiKey), porkbun.WithSecretKey(tc.secretKey))
			req, _ := http.NewRequest("GET", server.URL, nil)

			_, err := client.Do(req)
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
		})
	}
}
