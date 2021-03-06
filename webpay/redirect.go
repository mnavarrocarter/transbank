package webpay

import (
	"crypto/rand"
	_ "embed"
	"encoding/hex"
	"html/template"
	"io"
	"net/http"
)

//go:embed redirect.html
var redirectHtml string
var redirectTpl, _ = template.New("redirect").Parse(redirectHtml)

// Render renders the Webpay Payment Form using the passed http.Response writer
//
// It takes r as a source of bytes to generate a randomized form id to avoid token scraping.
//
// If r == nil then the rand.Reader is used
func (resp *CreateTransactionResponse) Render(rw http.ResponseWriter, r io.Reader) error {
	if r == nil {
		r = rand.Reader
	}

	b := make([]byte, 16)

	_, err := r.Read(b)
	if err != nil {
		return err
	}

	formId := hex.EncodeToString(b)

	rw.Header().Set("Content-Type", "text/html; charset=utf8")

	err = redirectTpl.Execute(rw, map[string]interface{}{
		"Url":    resp.Url,
		"Token":  resp.Token,
		"FormId": formId,
	})
	if err != nil {
		return err
	}
	return nil
}
