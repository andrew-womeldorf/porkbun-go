package porkbun

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

type Option func(*Client) error

type Client struct {
	apiKey    string
	secretKey string
	baseUrl   string
	client    HttpClient
}

// NewClient creates a new porkbun client.
// By default, it
func NewClient(options ...Option) (*Client, error) {
	c := &Client{
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
	return func(c *Client) error {
		c.apiKey = key
		return nil
	}
}

func WithSecretKey(key string) Option {
	return func(c *Client) error {
		c.secretKey = key
		return nil
	}
}

func WithBaseUrl(url string) Option {
	return func(c *Client) error {
		c.baseUrl = url
		return nil
	}
}

func WithHttpClient(httpClient HttpClient) Option {
	return func(c *Client) error {
		c.client = httpClient
		return nil
	}
}

func (c *Client) withAuthentication(body []byte) ([]byte, error) {
	if c.apiKey == "" {
		return nil, MissingAccessKeyError{Key: PORKBUN_API_KEY}
	}

	if c.secretKey == "" {
		return nil, MissingAccessKeyError{Key: PORKBUN_SECRET_KEY}
	}

	var orig map[string]interface{}

	if len(body) > 0 {
		if err := json.Unmarshal(body, &orig); err != nil {
			return nil, fmt.Errorf("could not unmarshal body, %w", err)
		}
	}

	newMap := map[string]interface{}{
		"apikey":       c.apiKey,
		"secretapikey": c.secretKey,
	}

	// Add original body
	for k, v := range orig {
		newMap[k] = v
	}

	// Marshal new body
	newBody, err := json.Marshal(newMap)
	if err != nil {
		return nil, fmt.Errorf("err marshaling json, %w", err)
	}

	return newBody, nil
}

func (c *Client) do(ctx context.Context, endpoint string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/%s", c.baseUrl, endpoint),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("err creating new request, %w", err)
	}
	return c.client.Do(req)
}
