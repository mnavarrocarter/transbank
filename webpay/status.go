package webpay

import (
	"context"
	"fmt"
	"net/http"
)

// Status returns the status of a transaction or error
//
// If the transaction was not successful, resp will be a pointer to TransactionInfo
// and err will be ErrTransaction explaining the cause of the error with detail. You
// should be careful to handle this particular case in your code.
//
// Any other error returned by this function is an unexpected error
func (c *Client) Status(ctx context.Context, token string) (*TransactionInfo, error) {
	resp := &TransactionInfo{}
	path := fmt.Sprintf("/rswebpaytransaction/api/webpay/v1.0/transactions/%s", token)

	err := c.sendRequest(ctx, http.MethodGet, path, nil, resp)
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != 0 {
		err = wrapTransactionError(resp.ResponseCode)
	}

	return resp, err
}
