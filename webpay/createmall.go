package webpay

import (
	"context"
	"net/http"
)

type MallTransaction struct {
	BuyOrder  string `json:"buy_order"`
	SessionId string `json:"session_id"`
	ReturnUrl string `json:"return_url"`
	Details   []struct {
		Amount       int    `json:"amount"`
		CommerceCode int64  `json:"commerce_code"`
		BuyOrder     string `json:"buy_order"`
	} `json:"details"`
}


func (c *Client) CreateMallTransaction(ctx context.Context, req *MallTransaction) (*CreateTransactionResponse, error) {
	resp := &CreateTransactionResponse{}
	err := c.sendRequest(ctx, http.MethodPost, "/rswebpaytransaction/api/webpay/v1.0/transactions", req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}