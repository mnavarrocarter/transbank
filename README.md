Golang Transbank SDK
====================

A solid and well tested Transbank SDK for Golang

```bash
go get github.com/mmavarrocarter/transbank
```

## Implemented Transbank Services

- [x] Webpay Plus - Normal
- [x] Webpay Plus - Deferred
- [ ] Webpay Plus - Mall
- [ ] OneClick - Mall
- [ ] Pos Integrado

## Webpay Usage

```go
package main

import (
	"context"
	"crypto/rand"
	"github.com/mnavarrocarter/transbank/webpay"
	"net/http"
)

func handleInitTransaction(client *webpay.Client) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		// In the real world, you obtain some of this data from another source
		transaction := &webpay.Transaction{
			BuyOrder:  "order1234",
			SessionId: "session1234",
			Amount: 12344,
			ReturnUrl: "https://my.app.com/webpay/callback",
		}

		// Create the transaction in webpay
		resp, err := client.CreateTransaction(context.Background(), transaction)

		// Handle any errors
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadGateway)
			return
		}

		// Render the webpay payment form with a self-posting js form
		err = resp.Render(rw, rand.Reader)
		if err != nil {
			panic(err)
        }
	}
}
```