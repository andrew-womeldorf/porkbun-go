package porkbun

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ApiError struct {
	Code int
	Body string
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.Code, e.Body)
}

type CreateDnsRecordRequest struct {
	// The subdomain for the record being created, not including the domain
	// itself. Leave blank to create a record on the root domain. Use * to
	// create a wildcard record.
	Name string `json:"name"`

	// The type of record being created. Valid types are: A, MX, CNAME, ALIAS,
	// TXT, NS, AAAA, SRV, TLSA, CAA.
	Type string `json:"type"`

	// The answer content for the record. Please see the DNS management popup
	// from the domain management console for proper formatting of each record
	// type.
	Content string `json:"content"`

	// The time to live in seconds for the record. The minimum and the default
	// is 600 seconds. Optional.
	TTL string `json:"ttl"`

	// The priority of the record for those that support it. Optional.
	Priority string `json:"prio"`
}

type CreateDnsRecordResponse struct {
	// A status indicating whether or not the command was successfuly
	// processed.
	Status string `json:"status"`

	// The ID of the record created.
	ID int `json:"id"`
}

// CreateDnsRecord creates a DNS entry in Porkbun.
//
// https://porkbun.com/api/json/v3/documentation#DNS%20Create%20Record
func (c *client) CreateDnsRecord(ctx context.Context, domain string, params *CreateDnsRecordRequest) (*CreateDnsRecordResponse, error) {
	reqBody, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("could not marshal params, %w", err)
	}

	body, err := c.withAuthentication(reqBody)
	if err != nil {
		return nil, fmt.Errorf("err adding authentication, %w", err)
	}

	res, err := c.do(ctx, fmt.Sprintf("/api/json/v3/dns/create/%s", domain), body)
	if err != nil {
		return nil, fmt.Errorf(
			"err creating dns record %q %q %q, %w",
			params.Name,
			params.Type,
			params.Content,
			err,
		)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		// Read response body
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		// Return custom error
		return nil, &ApiError{
			Code: res.StatusCode,
			Body: string(body),
		}
	}

	var response CreateDnsRecordResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&response); err != nil {
		return nil, fmt.Errorf("could not unmarshal response body, %w", err)
	}

	return &response, nil
}
