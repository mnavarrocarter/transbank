package webpay

import "time"

type RefundRequest struct {
	Amount int `json:"amount"`
}

type RefundResponse struct {
	Type              string    `json:"type"`
	AuthorizationCode string    `json:"authorization_code"`
	AuthorizationDate time.Time `json:"authorization_date"`
	NullifiedAmount   float64   `json:"nullified_amount"`
	Balance           float64   `json:"balance"`
	ResponseCode      int       `json:"response_code"`
}