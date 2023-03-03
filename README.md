# Movider API Client

[![Go](https://github.com/royge/movidersms/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/royge/movidersms/actions/workflows/go.yml)

## How to Use It

```go
package main

import (
	"context"
	"log"

	"github.com/royge/movidersms"
)

func main() {
	creds := movidersms.Credentials{"api-key", "api-secret"}
	sender := movidersms.NewSender(creds, []string{})

	res, err := sender.SendMessage(
		context.Background(),
		[]string{"6392612345678"},
		"Test Message",
	)
	if err != nil {
		panic(err)
	}

	log.Printf("Result: %+v\n", res)
}
```
