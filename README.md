# Smart Messaging Suite API Client

[![Go](https://github.com/royge/smartsms/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/royge/smartsms/actions/workflows/go.yml)

## How to Use It

```go
package main

import (
	"context"
	"log"

	"github.com/royge/smartsms"
)

func main() {
	creds := smartsms.Credentials{"test@email.com", "Secret"}
	sender := smartsms.NewSender(creds)

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
