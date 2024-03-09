# Porkbun Go Client

Manage DNS Records in [Porkbun][].

## Install

```sh
go get github.com/andrew-womeldorf/porkbun-go
```

## Usage

```go
package main

import "github.com/andrew-womeldorf/porkbun-go"

func main() {
    client := porkbun.NewClient(
        porkbun.WithApiKey("pk1_0000000000000000000000000000000000000000000000000000000000000000"),
        porkbun.WithSecretKey("sk1_0000000000000000000000000000000000000000000000000000000000000000"),
    )
}
```

Alternatively, you can set the access keys with environment variables:

- `PORKBUN_API_KEY`
- `PORKBUN_SECRET_KEY`

```go
client := porkbun.NewClient()
```

[porkbun]: https://porkbun.com
