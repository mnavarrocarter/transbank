package webpay

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type TransactionInfo struct {
	Vci                string      `json:"vci"`
	Amount             int         `json:"amount"`
	Status             string      `json:"status"`
	BuyOrder           string      `json:"buy_order"`
	SessionId          string      `json:"session_id"`
	CardDetail         *CardDetail `json:"card_detail,omitempty"`
	AccountingDate     string      `json:"accounting_date"`
	TransactionDate    time.Time   `json:"transaction_date"`
	AuthorizationCode  string      `json:"authorization_code"`
	PaymentTypeCode    string      `json:"payment_type_code"`
	ResponseCode       int         `json:"response_code"`
	InstallmentsNumber int         `json:"installments_number"`
}

type CardDetail struct {
	CardNumber string `json:"card_number,omitempty"`
}

// Confirm confirms a transaction in Transbank
//
// If the transaction was not successful, resp will be a pointer to TransactionInfo
// and err will be ErrTransaction explaining the cause of the error with detail according
// to the documentation found at https://www.transbankdevelopers.cl/producto/webpay#codigos-de-respuesta-de-autorizacion.
// You should be careful to handle this particular case in your code.
//
// Any other error returned by this function is an unexpected error
func (c *Client) Confirm(ctx context.Context, token string) (*TransactionInfo, error) {
	resp := &TransactionInfo{}
	path := fmt.Sprintf("/rswebpaytransaction/api/webpay/v1.0/transactions/%s", token)

	err := c.sendRequest(ctx, http.MethodPut, path, nil, resp)
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != 0 {
		err = wrapTransactionError(resp.ResponseCode)
	}

	return resp, err
}
