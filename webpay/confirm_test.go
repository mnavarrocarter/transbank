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

//go:embed testdata/confirm.0.res.json
var confirmRes0 []byte

//go:embed testdata/confirm.1.res.json
var confirmRes1 []byte

//go:embed testdata/confirm.2.res.json
var confirmRes2 []byte

func TestThatConfirmsTransaction(t *testing.T) {

	client := webpay.NewTestingClient()

	httptestclient.ExpectRequest(t).
		WithMethod("PUT").
		WithUrl("https://webpay3gint.transbank.cl/rswebpaytransaction/api/webpay/v1.0/transactions/01abbc6652ded78ae8c0a9dd50e4e4b2bdc977f8832f5771d542a0734131ef51").
		WithHeaders(reqHeaders).
		WillReturnResponse().
		WithStatusCode(200).
		WithBody(confirmRes0).
		WithHeaders(resHeaders).
		InjectInDefaultClient()

	defer httptestclient.RestoreTransport()

	resp, err := client.Confirm(context.Background(), "01abbc6652ded78ae8c0a9dd50e4e4b2bdc977f8832f5771d542a0734131ef51")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "ordenCompra12345678", resp.BuyOrder)
	assert.Equal(t, "sesion1234557545", resp.SessionId)
	assert.Equal(t, "6623", resp.CardDetail.CardNumber)
}

func TestConfirmTransactionWithTransactionError(t *testing.T) {

	client := webpay.NewTestingClient()

	httptestclient.ExpectRequest(t).
		WithMethod("PUT").
		WithUrl("https://webpay3gint.transbank.cl/rswebpaytransaction/api/webpay/v1.0/transactions/01abbc6652ded78ae8c0a9dd50e4e4b2bdc977f8832f5771d542a0734131ef51").
		WithHeaders(reqHeaders).
		WillReturnResponse().
		WithStatusCode(200).
		WithBody(confirmRes1).
		WithHeaders(resHeaders).
		InjectInDefaultClient()

	defer httptestclient.RestoreTransport()

	resp, err := client.Confirm(context.Background(), "01abbc6652ded78ae8c0a9dd50e4e4b2bdc977f8832f5771d542a0734131ef51")

	assert.Error(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "transaction error: generic transaction error", err.Error())
	assert.Equal(t, "ordenCompra12345678", resp.BuyOrder)
	assert.Equal(t, "sesion1234557545", resp.SessionId)
	assert.Equal(t, "6623", resp.CardDetail.CardNumber)
	assert.Equal(t, -3, resp.ResponseCode)
}

func TestConfirmWithInvalidToken(t *testing.T) {

	client := webpay.NewTestingClient()

	httptestclient.ExpectRequest(t).
		WithMethod("PUT").
		WithUrl("https://webpay3gint.transbank.cl/rswebpaytransaction/api/webpay/v1.0/transactions/3265236523623965723865233324").
		WithHeaders(reqHeaders).
		WillReturnResponse().
		WithStatusCode(422).
		WithBody(confirmRes2).
		WithHeaders(resHeaders).
		InjectInDefaultClient()

	defer httptestclient.RestoreTransport()

	resp, err := client.Confirm(context.Background(), "3265236523623965723865233324")

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestConfirmHandlesConnectionError(t *testing.T) {
	client := webpay.NewTestingClient()

	httptestclient.ExpectRequest(t).
		WithMethod("POST").
		WithUrl("https://webpay3gint.transbank.cl/rswebpaytransaction/api/webpay/v1.0/transactions/01abbc6652ded78ae8c0a9dd50e4e4b2bdc977f8832f5771d542a0734131ef51").
		WithHeaders(reqHeaders).
		WillReturnResponse().
		WithError(errors.New("connection error")).
		InjectInDefaultClient()

	defer httptestclient.RestoreTransport()

	resp, err := client.Confirm(context.Background(), "01abbc6652ded78ae8c0a9dd50e4e4b2bdc977f8832f5771d542a0734131ef51")

	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.Equal(t, "unexpected error sending request: Put \"https://webpay3gint.transbank.cl/rswebpaytransaction/api/webpay/v1.0/transactions/01abbc6652ded78ae8c0a9dd50e4e4b2bdc977f8832f5771d542a0734131ef51\": connection error", err.Error())
}
