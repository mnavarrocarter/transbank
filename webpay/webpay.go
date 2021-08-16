package webpay

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

const (
	prodUrl = "https://webpay3g.transbank.cl"
	intUrl = "https://webpay3gint.transbank.cl"
)

const (
	// intApiToken is the api token for the testing environment
	intApiToken = "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C"

	IntNormalCommerceCode = "597055555532"
	IntDeferredCommerceCode = "597055555540"
	IntMallCommerceCode = "597055555535"
	IntMallChildOneCommerceCode = "597055555536"
	IntMallChildTwoCommerceCode = "597055555537"
	IntMallDeferredCommerceCode = "597055555581"
	IntMallDeferredChildOneCommerceCode = "597055555582"
	IntMallDeferredChildTwoCommerceCode = "597055555583"
)

// A Client performs requests to the Webpay REST api
type Client struct {
	apiToken string
	commerceCode string
	baseUrl string
	client *http.Client
}

// NewClient creates a Webpay client
// If client = nil, then a normal http.Client instance is created
func NewClient(apiToken, commerceCode string, production bool, client *http.Client) *Client {
	var baseUrl = intUrl
	if production == true {
		baseUrl = prodUrl
	}
	if client == nil {
		client = &http.Client{}
	}
	return &Client{apiToken: apiToken, commerceCode: commerceCode, baseUrl: baseUrl, client: client}
}

// NewTestingClient creates a client already configured for the testing environment
func NewTestingClient() *Client  {
	return NewClient(intApiToken, IntNormalCommerceCode, false, nil)
}

// SetCommerceCode allows to override the main commerce code
func (c * Client) SetCommerceCode(code string)  {
	c.commerceCode = code
}

// SetHttpClient swaps the inner http.Client instance this Client uses
// This is very useful for testing purposes
func (c * Client) SetHttpClient(client *http.Client)  {
	c.client = client
}

func (c *Client) sendRequest(method, path string, input, output interface{}) error {
	b, err := json.Marshal(input)
	if err != nil {
		ErrUnexpected.Err = err
		return ErrUnexpected
	}

	buff := bytes.NewBuffer(b)

	req, err := http.NewRequest(method, c.baseUrl + path, buff)

	if err != nil {
		ErrUnexpected.Err = err
		return ErrUnexpected
	}

	req.Header.Set("Tbk-Api-Key-Id", c.commerceCode)
	req.Header.Set("Tbk-Api-Key-Secret", c.apiToken)
	req.Header.Set("Content-Type", "application/json; charset=utf8")
	req.Header.Set("Accept", "application/json; charset=utf8")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		ErrHttpProtocol.Response = resp
		ErrHttpProtocol.StatusCode = resp.StatusCode
		return ErrHttpProtocol
	}

	err = json.NewDecoder(resp.Body).Decode(output)
	if err != nil {
		ErrUnexpected.Err = err
		return ErrUnexpected
	}

	return nil
}