package webpay_test

import (
	"context"
	_ "embed"
	"errors"
	"github.com/mnavarrocarter/httpclientmock"
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

var normalTransaction = &webpay.Transaction{
	BuyOrder:  "ordenCompra12345678",
	SessionId: "sesion1234557545",
	Amount:    10000,
	ReturnUrl: "http://www.comercio.cl/webpay/retorno",
}

var createTests = []struct {
	name        string
	transaction *webpay.Transaction
	mock        *httpclientmock.Mock
	assertions  func(t *testing.T, resp *webpay.CreateTransactionResponse, err error)
}{
	{
		name:        "it creates successfully",
		transaction: normalTransaction,
		mock: &httpclientmock.Mock{
			Expect: &httpclientmock.Request{
				Method:  "POST",
				Url:     "https://webpay3gint.transbank.cl/rswebpaytransaction/api/webpay/v1.0/transactions",
				Headers: reqHeaders,
				Body:    createReq0,
			},
			Return: &httpclientmock.Response{
				StatusCode: 200,
				Headers:    resHeaders,
				Body:       createRes0,
			},
		},
		assertions: func(t *testing.T, resp *webpay.CreateTransactionResponse, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, "e9d555262db0f989e49d724b4db0b0af367cc415cde41f500a776550fc5fddd3", resp.Token)
			assert.Equal(t, "https://webpay3gint.transbank.cl/webpayserver/initTransaction", resp.Url)
		},
	},
	{
		name: "it handles bad input error successfully",
		transaction: &webpay.Transaction{
			BuyOrder:  "ordenCompra12345678",
			SessionId: "sesion1234557545",
			Amount:    10000,
		},
		mock: &httpclientmock.Mock{
			Expect: &httpclientmock.Request{
				Method:  "POST",
				Url:     "https://webpay3gint.transbank.cl/rswebpaytransaction/api/webpay/v1.0/transactions",
				Headers: reqHeaders,
				Body:    createReq1,
			},
			Return: &httpclientmock.Response{
				StatusCode: 422,
				Headers:    resHeaders,
				Body:       createRes1,
			},
		},
		assertions: func(t *testing.T, resp *webpay.CreateTransactionResponse, err error) {
			assert.Nil(t, resp)
			assert.NotNil(t, err)
			assert.Equal(t, "request returned status 422 saying: return_url is required!", err.Error())
		},
	},
	{
		name:        "it handles connection error correctly",
		transaction: normalTransaction,
		mock: &httpclientmock.Mock{
			Expect: &httpclientmock.Request{
				Method:  "POST",
				Url:     "https://webpay3gint.transbank.cl/rswebpaytransaction/api/webpay/v1.0/transactions",
				Headers: reqHeaders,
				Body:    createReq0,
			},
			Return: &httpclientmock.Response{
				Err: errors.New("connection error"),
			},
		},
		assertions: func(t *testing.T, resp *webpay.CreateTransactionResponse, err error) {
			assert.Nil(t, resp)
			assert.NotNil(t, err)
			assert.Equal(t, "unexpected error sending request: Post \"https://webpay3gint.transbank.cl/rswebpaytransaction/api/webpay/v1.0/transactions\": connection error", err.Error())
		},
	},
}

func TestCreate(t *testing.T) {
	client := webpay.NewTestingClient()

	for _, test := range createTests {
		t.Run(test.name, func(t *testing.T) {
			restore := test.mock.InjectInClient(t, nil)
			defer restore()
			resp, err := client.CreateTransaction(context.Background(), test.transaction)
			test.assertions(t, resp, err)
		})
	}
}
