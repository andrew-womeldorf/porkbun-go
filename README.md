# Porkbun Go Client

A relatively simple wrapper around the [porkbun api
v3](https://porkbun.com/api/json/v3/documentation). Currently only supports DNS
management.

## Install

```sh
go install github.com/andrew-womeldorf/porkbun-go/cmd/porkbun@latest
```

## Authentication

All API calls must include valid API keys. You can create API keys at
[porkbun.com/account/api](porkbun.com/account/api). You can test communication
with the API using the `ping` endpoint. The ping endpoint will also return your
IP address, this can be handy when building dynamic DNS clients.

You can optionally set the credentials with the following environment
variables:

- `PORKBUN_API_KEY`
- `PORKBUN_SECRET_KEY`

**Important:** To manage a domain's DNS via the API, you must toggle the `API
ACCESS` setting within the management console for each domain you want to
manage programmatically.

## Usage

**As a cli tool:** `porkbun help`.

**As a library:**

See
[pkg.go.dev/github.com/andrew-womeldorf/porkbun-go](https://pkg.go.dev/github.com/andrew-womeldorf/porkbun-go)
for library docs.

```go
package main

import (
	"fmt"

	"github.com/andrew-womeldorf/porkbun-go"
)

func main() {
	client := porkbun.NewClient(
		porkbun.WithApiKey("pk1_0000000000000000000000000000000000000000000000000000000000000000"),
		porkbun.WithSecretKey("sk1_0000000000000000000000000000000000000000000000000000000000000000"),
	)

	// alternatively, if environment variables are set:
	// client := porkbun.NewClient()

	res, err := client.Ping(ctx)
	if err != nil {
		t.Fatalf("got %s, want nil", err)
	}

	fmt.Println(res.Status, res.YourIP)
}
```
