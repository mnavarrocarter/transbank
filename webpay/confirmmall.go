package webpay

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type MallTransactionInfo struct {
	BuyOrder   string `json:"buy_order"`
	CardDetail struct {
		CardNumber string `json:"card_number"`
	} `json:"card_detail"`
	AccountingDate  string    `json:"accounting_date"`
	TransactionDate time.Time `json:"transaction_date"`
	Vci             string    `json:"vci"`
	Details         []struct {
		Amount             int    `json:"amount"`
		Status             string `json:"status"`
		AuthorizationCode  string `json:"authorization_code"`
		PaymentTypeCode    string `json:"payment_type_code"`
		ResponseCode       int    `json:"response_code"`
		InstallmentsNumber int    `json:"installments_number"`
		CommerceCode       string `json:"commerce_code"`
		BuyOrder           string `json:"buy_order"`
	} `json:"details"`
}

func (c *Client) ConfirmMall(ctx context.Context, token string) (*MallTransactionInfo, error) {
	path := fmt.Sprintf("/rswebpaytransaction/api/webpay/v1.0/transactions/%s", token)
	resp := &MallTransactionInfo{}
	err := c.sendRequest(ctx, http.MethodPut, path, nil, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}