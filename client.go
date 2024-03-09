package porkbun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	PORKBUN_API_KEY    = "PORKBUN_API_KEY"
	PORKBUN_SECRET_KEY = "PORKBUN_SECRET_KEY"
)

type MissingAccessKeyError struct {
	Key string
}

func (e MissingAccessKeyError) Error() string {
	var keyType string

	if e.Key == PORKBUN_API_KEY {
		keyType = "api"
	}

	if e.Key == PORKBUN_SECRET_KEY {
		keyType = "secret"
	}

	return fmt.Sprintf("missing porkbun %q key. try setting %q to the environment", keyType, e.Key)
}

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Option func(*client) error

type client struct {
	apiKey    string
	secretKey string
	baseUrl   string
	client    HttpClient
}

// NewClient creates a new porkbun client.
// By default, it
func NewClient(options ...Option) (*client, error) {
	c := &client{
		apiKey:    os.Getenv(PORKBUN_API_KEY),
		secretKey: os.Getenv(PORKBUN_SECRET_KEY),
		baseUrl:   "https://porkbun.com",
		client:    &http.Client{},
	}

	for _, option := range options {
		err := option(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func WithApiKey(key string) Option {
	return func(c *client) error {
		c.apiKey = key
		return nil
	}
}

func WithSecretKey(key string) Option {
	return func(c *client) error {
		c.secretKey = key
		return nil
	}
}

func WithBaseUrl(url string) Option {
	return func(c *client) error {
		c.baseUrl = url
		return nil
	}
}

func WithHttpClient(httpClient HttpClient) Option {
	return func(c *client) error {
		c.client = httpClient
		return nil
	}
}

// Do sends an HTTP request and return an HTTP response.
// It automatically sends the APIKey and SecretAPIKey in the request body on
// your behalf.
func (c *client) Do(req *http.Request) (*http.Response, error) {
	if c.apiKey == "" {
		return nil, MissingAccessKeyError{Key: PORKBUN_API_KEY}
	}

	if c.secretKey == "" {
		return nil, MissingAccessKeyError{Key: PORKBUN_SECRET_KEY}
	}

	var orig map[string]interface{}

	if req.Body != nil {
		defer req.Body.Close()

		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&orig); err != nil {
			return nil, fmt.Errorf("could not unmarshal body, %w", err)
		}
	}

	newMap := map[string]interface{}{
		"apiKey":    c.apiKey,
		"secretKey": c.secretKey,
	}

	// Add original body
	for k, v := range orig {
		newMap[k] = v
	}

	// Marshal new body
	newBody, err := json.Marshal(newMap)
	if err != nil {
		return nil, err
	}

	// Set new request body
	req.Body = io.NopCloser(bytes.NewReader(newBody))
	req.ContentLength = int64(len(newBody))

	return c.client.Do(req)
}
