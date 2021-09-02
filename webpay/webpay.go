// Package webpay implements a Client for accessing the Webpay API in a programmatic way
package webpay

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
)

const (
	ProdUrl = "https://webpay3g.transbank.cl"
	IntUrl  = "https://webpay3gint.transbank.cl"
)

const (
	// IntApiToken is the api token for the testing environment
	IntApiToken = "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C"

	IntNormalCommerceCode               = "597055555532"
	IntDeferredCommerceCode             = "597055555540"
	IntMallCommerceCode                 = "597055555535"
	IntMallChildOneCommerceCode         = "597055555536"
	IntMallChildTwoCommerceCode         = "597055555537"
	IntMallDeferredCommerceCode         = "597055555581"
	IntMallDeferredChildOneCommerceCode = "597055555582"
	IntMallDeferredChildTwoCommerceCode = "597055555583"
)

// A Client performs requests to the Webpay REST api
type Client struct {
	apiToken     string
	commerceCode string
	baseUrl      string
	client       *http.Client
}

// NewClient creates a Webpay client
// If client = nil, then http.DefaultClient is used
func NewClient(apiToken, commerceCode string, production bool, client *http.Client) *Client {
	var baseUrl = IntUrl
	if production == true {
		baseUrl = ProdUrl
	}
	if client == nil {
		client = http.DefaultClient
	}
	return &Client{apiToken: apiToken, commerceCode: commerceCode, baseUrl: baseUrl, client: client}
}

// NewTestingClient creates a client already configured for the testing environment
func NewTestingClient() *Client {
	return NewClient(IntApiToken, IntNormalCommerceCode, false, nil)
}

// SetCommerceCode allows to override the main commerce code
func (c *Client) SetCommerceCode(code string) {
	c.commerceCode = code
}

// SetHttpClient swaps the inner http.Client instance this Client use
//
// This is very useful for testing purposes
func (c *Client) SetHttpClient(client *http.Client) {
	c.client = client
}

func (c *Client) sendRequest(ctx context.Context, method, path string, input, output interface{}) error {
	b, err := json.MarshalIndent(input, "", "    ")
	if err != nil {
		return errors.Wrapf(err, "unexpected error serializing json payload")
	}

	buff := bytes.NewBuffer(b)

	req, err := http.NewRequest(method, c.baseUrl+path, buff)
	if err != nil {
		return errors.Wrapf(err, "unexpected error creating request")
	}

	req = req.WithContext(ctx)

	req.Header.Set("Tbk-Api-Key-Id", c.commerceCode)
	req.Header.Set("Tbk-Api-Key-Secret", c.apiToken)
	req.Header.Set("Content-Type", "application/json; charset=utf8")
	req.Header.Set("Accept", "application/json; charset=utf8")

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "unexpected error sending request")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("could not close file: %s\n", err.Error())
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var errmsg map[string]string
		err := json.NewDecoder(resp.Body).Decode(&errmsg)
		if err != nil {
			return errors.Wrapf(err, "unexpected error decoding json response")
		}
		return errors.Errorf("request returned status %d saying: %s", resp.StatusCode, errmsg["error_message"])
	}

	err = json.NewDecoder(resp.Body).Decode(output)
	if err != nil {
		return errors.Wrapf(err, "unexpected error decoding json response")
	}

	return nil
}
