package porkbun

import (
	"context"
	"encoding/json"
	"fmt"
)

type PingResponse struct {
	Status string `json:"status"`
	YourIP string `json:"yourIp"`
}

// Ping tests communication with the API using the ping endpoint. The ping
// endpoint will also return your IP address, this can be handy when building
// dynamic DNS clients.
func (c *Client) Ping(ctx context.Context) (*PingResponse, error) {
	body, err := c.withAuthentication(nil)
	if err != nil {
		return nil, fmt.Errorf("err adding authentication, %w", err)
	}

	res, err := c.do(ctx, "/api/json/v3/ping", body)
	if err != nil {
		return nil, fmt.Errorf("err calling ping, %w", err)
	}
	defer res.Body.Close()

	var response PingResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&response); err != nil {
		return nil, fmt.Errorf("could not unmarshal response body, %w", err)
	}

	return &response, nil
}
