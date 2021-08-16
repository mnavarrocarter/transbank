package webpay

import (
	"fmt"
	"net/http"
	"time"
)

type TransactionInfo struct {
	Vci        string `json:"vci"`
	Amount     int    `json:"amount"`
	Status     string `json:"status"`
	BuyOrder   string `json:"buy_order"`
	SessionId  string `json:"session_id"`
	CardDetail struct {
		CardNumber string `json:"card_number,omitempty"`
	} `json:"card_detail"`
	AccountingDate     string    `json:"accounting_date"`
	TransactionDate    time.Time `json:"transaction_date"`
	AuthorizationCode  string    `json:"authorization_code"`
	PaymentTypeCode    string    `json:"payment_type_code"`
	ResponseCode       int       `json:"response_code"`
	InstallmentsNumber int       `json:"installments_number"`
}

// Confirm confirms a transaction in Transbank
func (c *Client) Confirm(token string) (resp *TransactionInfo, err error) {
	path := fmt.Sprintf("/rswebpaytransaction/api/webpay/v1.0/transactions/%s", token)

	err = c.sendRequest(http.MethodPut, path, nil, resp)
	if err != nil {
		return nil, err
	}

	err = mapError(resp.ResponseCode)
	if err != nil {
		return resp, err
	}

	return resp, nil
}