package webpay_test

import (
	"context"
	_ "embed"
	"github.com/mnavarrocarter/httptestclient"
	"github.com/pkg/errors"
	"testing"

	"github.com/mnavarrocarter/transbank/webpay"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/create.0.req.json
var createReq0 []byte

//go:embed testdata/create.1.req.json
var createReq1 []byte

//go:embed testdata/create.0.res.json
var createRes0 []byte

//go:embed testdata/create.1.res.json
var createRes1 []byte

func TestThatCreatesTransaction(t *testing.T) {

	client := webpay.NewTestingClient()

	httptestclient.ExpectRequest(t).
		WithMethod("POST").
		WithUrl("https://webpay3gint.transbank.cl/rswebpaytransaction/api/webpay/v1.0/transactions").
		WithBody(createReq0).
		WithHeaders(reqHeaders).
		WillReturnResponse().
		WithStatusCode(200).
		WithBody(createRes0).
		WithHeaders(resHeaders).
		InjectInDefaultClient()

	defer httptestclient.RestoreTransport()

	resp, err := client.CreateTransaction(context.Background(), &webpay.Transaction{
		BuyOrder:  "ordenCompra12345678",
		SessionId: "sesion1234557545",
		Amount:    10000,
		ReturnUrl: "http://www.comercio.cl/webpay/retorno",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "e9d555262db0f989e49d724b4db0b0af367cc415cde41f500a776550fc5fddd3", resp.Token)
	assert.Equal(t, "https://webpay3gint.transbank.cl/webpayserver/initTransaction", resp.Url)
}

func TestThatHandlesBadInput(t *testing.T) {
	client := webpay.NewTestingClient()

	httptestclient.ExpectRequest(t).
		WithMethod("POST").
		WithUrl("https://webpay3gint.transbank.cl/rswebpaytransaction/api/webpay/v1.0/transactions").
		WithBody(createReq1).
		WithHeaders(reqHeaders).
		WillReturnResponse().
		WithStatusCode(422).
		WithBody(createRes1).
		WithHeaders(resHeaders).
		InjectInDefaultClient()

	defer httptestclient.RestoreTransport()

	resp, err := client.CreateTransaction(context.Background(), &webpay.Transaction{
		BuyOrder:  "ordenCompra12345678",
		SessionId: "sesion1234557545",
		Amount:    10000,
		ReturnUrl: "",
	})

	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.Equal(t, "request returned status 422 saying: return_url is required!", err.Error())
}

func TestThatHandlesConnectionError(t *testing.T) {
	client := webpay.NewTestingClient()

	httptestclient.ExpectRequest(t).
		WithMethod("POST").
		WithUrl("https://webpay3gint.transbank.cl/rswebpaytransaction/api/webpay/v1.0/transactions").
		WithBody(createReq0).
		WithHeaders(reqHeaders).
		WillReturnResponse().
		WithError(errors.New("connection error")).
		InjectInDefaultClient()

	defer httptestclient.RestoreTransport()

	resp, err := client.CreateTransaction(context.Background(), &webpay.Transaction{
		BuyOrder:  "ordenCompra12345678",
		SessionId: "sesion1234557545",
		Amount:    10000,
		ReturnUrl: "http://www.comercio.cl/webpay/retorno",
	})

	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.Equal(t, "unexpected error sending request: Post \"https://webpay3gint.transbank.cl/rswebpaytransaction/api/webpay/v1.0/transactions\": connection error", err.Error())
}
