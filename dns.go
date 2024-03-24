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

type Record struct {
	Id    string `json:"id"`
	Notes string `json:"notes"`

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

	// The Id of the record created.
	//
	// Creating a record, the Id returned is an int, but every other method
	// expects Id to be a string. This is accomodating for the upstream API...
	Id int `json:"id"`
}

type DnsRecordsResponse struct {
	Status  string   `json:"status"`
	Records []Record `json:"records"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

// CreateDnsRecord creates a DNS entry in Porkbun.
//
// https://porkbun.com/api/json/v3/documentation#DNS%20Create%20Record
func (c *Client) CreateDnsRecord(ctx context.Context, domain string, params *Record) (*CreateDnsRecordResponse, error) {
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

// ListDnsRecords returns a list of DNS records.
// Get all available records by leaving the subdomain and recordType as empty.
// Find a subset of records by providing the subdomain and type.
func (c *Client) ListDnsRecords(ctx context.Context, domain, subdomain, recordType string) (*DnsRecordsResponse, error) {
	body, err := c.withAuthentication(nil)
	if err != nil {
		return nil, fmt.Errorf("err adding authentication, %w", err)
	}

	var url string
	if recordType != "" {
		url = fmt.Sprintf("/api/json/v3/dns/retrieveByNameType/%s/%s/%s", domain, recordType, subdomain)
	} else {
		url = fmt.Sprintf("/api/json/v3/dns/retrieve/%s", domain)
	}

	res, err := c.do(ctx, url, body)
	if err != nil {
		return nil, fmt.Errorf("err retrieving dns records, %w", err)
	}
	defer res.Body.Close()

	var response DnsRecordsResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&response); err != nil {
		return nil, fmt.Errorf("could not unmarshal response body, %w", err)
	}

	return &response, nil
}

func (c *Client) GetDnsRecordById(ctx context.Context, domain string, id int) (*DnsRecordsResponse, error) {
	body, err := c.withAuthentication(nil)
	if err != nil {
		return nil, fmt.Errorf("err adding authentication, %w", err)
	}

	res, err := c.do(ctx, fmt.Sprintf("/api/json/v3/dns/retrieve/%s/%d", domain, id), body)
	if err != nil {
		return nil, fmt.Errorf("err retrieving dns record, %w", err)
	}
	defer res.Body.Close()

	var response DnsRecordsResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&response); err != nil {
		return nil, fmt.Errorf("could not unmarshal response body, %w", err)
	}

	return &response, nil
}

// ModifyDnsRecord changes a DNS entry in Porkbun.
//
// Only Content, TTL, and Priority are necessary fields on the record.
//
// If record.Id is not empty, then modify a record found by the provided ID.
// Otherwise, the record will be looked up by the subdomain and type.
//
// https://porkbun.com/api/json/v3/documentation#DNS%20Edit%20Record%20by%20Domain%20and%20ID
func (c *Client) ModifyDnsRecord(ctx context.Context, domain string, record *Record) (*StatusResponse, error) {
	reqBody, err := json.Marshal(record)
	if err != nil {
		return nil, fmt.Errorf("could not marshal record, %w", err)
	}

	body, err := c.withAuthentication(reqBody)
	if err != nil {
		return nil, fmt.Errorf("err adding authentication, %w", err)
	}

	var url string
	if record.Id != "" {
		url = fmt.Sprintf("/api/json/v3/dns/edit/%s/%s", domain, record.Id)
	} else {
		if record.Type == "" {
			return nil, fmt.Errorf("record.Type must be set to modify this entry")
		}
		url = fmt.Sprintf("/api/json/v3/dns/editByNameType/%s/%s/%s", domain, record.Type, record.Name)
	}

	res, err := c.do(ctx, url, body)
	if err != nil {
		return nil, fmt.Errorf(
			"err editing dns record %q %q %q %q, %w",
			record.Id,
			record.Name,
			record.Type,
			record.Content,
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

	var response StatusResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&response); err != nil {
		return nil, fmt.Errorf("could not unmarshal response body, %w", err)
	}

	return &response, nil
}

// DeleteDnsRecordById deletes a DNS entry in Porkbun, looking up by id.
//
// https://porkbun.com/api/json/v3/documentation#DNS%20Delete%20Record%20by%20Domain%20and%20ID
func (c *Client) DeleteDnsRecordById(ctx context.Context, domain, id string) (*StatusResponse, error) {
	body, err := c.withAuthentication(nil)
	if err != nil {
		return nil, fmt.Errorf("err adding authentication, %w", err)
	}

	res, err := c.do(ctx, fmt.Sprintf("/api/json/v3/dns/delete/%s/%s", domain, id), body)
	if err != nil {
		return nil, fmt.Errorf("err deleting dns record %q, %w", id, err)
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

	var response StatusResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&response); err != nil {
		return nil, fmt.Errorf("could not unmarshal response body, %w", err)
	}

	return &response, nil
}

// DeleteDnsRecordByLookup deletes a DNS entry in Porkbun, looking up by subdomain and record type.
//
// https://porkbun.com/api/json/v3/documentation#DNS%20Delete%20Records%20by%20Domain,%20Subdomain%20and%20Type
func (c *Client) DeleteDnsRecordByLookup(ctx context.Context, domain, subdomain, recordType string) (*StatusResponse, error) {
	body, err := c.withAuthentication(nil)
	if err != nil {
		return nil, fmt.Errorf("err adding authentication, %w", err)
	}

	res, err := c.do(ctx, fmt.Sprintf("/api/json/v3/dns/deleteByNameType/%s/%s/%s", domain, recordType, subdomain), body)
	if err != nil {
		return nil, fmt.Errorf("err deleting dns record %q, %q, %w", subdomain, recordType, err)
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

	var response StatusResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&response); err != nil {
		return nil, fmt.Errorf("could not unmarshal response body, %w", err)
	}

	return &response, nil
}
