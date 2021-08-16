package webpay

import "time"

type CaptureRequest struct {
	BuyOrder          string `json:"buy_order"`
	AuthorizationCode string `json:"authorization_code"`
	CaptureAmount     int    `json:"capture_amount"`
}

type CaptureResponse struct {
	Token             string    `json:"token"`
	AuthorizationCode string    `json:"authorization_code"`
	AuthorizationDate time.Time `json:"authorization_date"`
	CapturedAmount    int       `json:"captured_amount"`
	ResponseCode      int       `json:"response_code"`
}