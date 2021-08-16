package webpay_test

import (
	_ "embed"
	"github.com/mnavarrocarter/transbank/webpay"
	"github.com/stretchr/testify/assert"
	"io"
	"math/rand"
	"net/http/httptest"
	"testing"
)

//go:embed redirect_test.html
var resultHtml string

func TestRender(t *testing.T)  {
	rw := httptest.NewRecorder()
	r := rand.New(rand.NewSource(0))

	tr := &webpay.CreateTransactionResponse{
		Token: "3257320957329",
		Url:   "https://example.com",
	}

	tr.Render(rw, r)

	resp := rw.Result()
	defer resp.Body.Close()

	html, _ := io.ReadAll(resp.Body)
	assert.Equal(t, resultHtml, string(html))
}