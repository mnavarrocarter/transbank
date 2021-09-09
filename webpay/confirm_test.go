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

//go:embed testdata/confirm.0.res.json
var confirmRes0 []byte

//go:embed testdata/confirm.1.res.json
var confirmRes1 []byte

//go:embed testdata/confirm.2.res.json
var confirmRes2 []byte

var confirmTests = []struct {
	name       string
	token      string
	mock       *httpclientmock.Mock
	assertions func(t *testing.T, resp *webpay.TransactionInfo, err error)
}{
	{
		name:  "it confirms successfully",
		token: "01abbc6652ded78ae8c0a9dd50e4e4b2bdc977f8832f5771d542a0734131ef51",
		mock: &httpclientmock.Mock{
			Expect: &httpclientmock.Request{
				Method:  "PUT",
				Url:     "https://webpay3gint.transbank.cl/rswebpaytransaction/api/webpay/v1.0/transactions/01abbc6652ded78ae8c0a9dd50e4e4b2bdc977f8832f5771d542a0734131ef51",
				Headers: reqHeaders,
			},
			Return: &httpclientmock.Response{
				StatusCode: 200,
				Headers:    resHeaders,
				Body:       confirmRes0,
			},
		},
		assertions: func(t *testing.T, resp *webpay.TransactionInfo, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, "ordenCompra12345678", resp.BuyOrder)
			assert.Equal(t, "sesion1234557545", resp.SessionId)
			assert.Equal(t, "6623", resp.CardDetail.CardNumber)
		},
	},
	{
		name:  "it handles transaction error",
		token: "01abbc6652ded78ae8c0a9dd50e4e4b2bdc977f8832f5771d542a0734131ef51",
		mock: &httpclientmock.Mock{
			Expect: &httpclientmock.Request{
				Method:  "PUT",
				Url:     "https://webpay3gint.transbank.cl/rswebpaytransaction/api/webpay/v1.0/transactions/01abbc6652ded78ae8c0a9dd50e4e4b2bdc977f8832f5771d542a0734131ef51",
				Headers: reqHeaders,
			},
			Return: &httpclientmock.Response{
				StatusCode: 200,
				Headers:    resHeaders,
				Body:       confirmRes1,
			},
		},
		assertions: func(t *testing.T, resp *webpay.TransactionInfo, err error) {
			assert.Error(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, "transaction error: generic transaction error", err.Error())
			assert.Equal(t, "ordenCompra12345678", resp.BuyOrder)
			assert.Equal(t, "sesion1234557545", resp.SessionId)
			assert.Equal(t, "6623", resp.CardDetail.CardNumber)
			assert.Equal(t, -3, resp.ResponseCode)
		},
	},
	{
		name:  "it handles invalid token",
		token: "3265236523623965723865233324",
		mock: &httpclientmock.Mock{
			Expect: &httpclientmock.Request{
				Method:  "PUT",
				Url:     "https://webpay3gint.transbank.cl/rswebpaytransaction/api/webpay/v1.0/transactions/3265236523623965723865233324",
				Headers: reqHeaders,
			},
			Return: &httpclientmock.Response{
				StatusCode: 422,
				Headers:    resHeaders,
				Body:       confirmRes2,
			},
		},
		assertions: func(t *testing.T, resp *webpay.TransactionInfo, err error) {
			assert.Error(t, err)
			assert.Nil(t, resp)
		},
	},
	{
		name:  "it handles connection error",
		token: "01abbc6652ded78ae8c0a9dd50e4e4b2bdc977f8832f5771d542a0734131ef51",
		mock: &httpclientmock.Mock{
			Expect: &httpclientmock.Request{
				Method:  "PUT",
				Url:     "https://webpay3gint.transbank.cl/rswebpaytransaction/api/webpay/v1.0/transactions/01abbc6652ded78ae8c0a9dd50e4e4b2bdc977f8832f5771d542a0734131ef51",
				Headers: reqHeaders,
			},
			Return: &httpclientmock.Response{
				Err: errors.New("connection error"),
			},
		},
		assertions: func(t *testing.T, resp *webpay.TransactionInfo, err error) {
			assert.Nil(t, resp)
			assert.NotNil(t, err)
			assert.Equal(t, "unexpected error sending request: Put \"https://webpay3gint.transbank.cl/rswebpaytransaction/api/webpay/v1.0/transactions/01abbc6652ded78ae8c0a9dd50e4e4b2bdc977f8832f5771d542a0734131ef51\": connection error", err.Error())
		},
	},
}

func TestConfirm(t *testing.T) {
	client := webpay.NewTestingClient()

	for _, test := range confirmTests {
		t.Run(test.name, func(t *testing.T) {
			restore := test.mock.InjectInClient(t, nil)
			defer restore()
			resp, err := client.Confirm(context.Background(), test.token)
			test.assertions(t, resp, err)
		})
	}
}
