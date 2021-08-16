package webpay

import (
	_ "embed"
	"net/http"
)

type Transaction struct {
	BuyOrder  string `json:"buy_order"`
	SessionId string `json:"session_id,omitempty"`
	Amount    int    `json:"amount"`
	ReturnUrl string `json:"return_url"`
}

type CreateTransactionResponse struct {
	Token string `json:"token"`
	Url   string `json:"url"`
}

// CreateTransaction creates a webpay normal transaction
func (c *Client) CreateTransaction(req *Transaction) (*CreateTransactionResponse, error) {
	resp := &CreateTransactionResponse{}
	err := c.sendRequest(http.MethodPost, "/rswebpaytransaction/api/webpay/v1.0/transactions", req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}