package webpay

import (
	"errors"
	"fmt"
	"net/http"
)

// When UseErrorLevelTwo is set to true, it maps the error codes to the new error definitions added in 2021.
// See https://www.transbankdevelopers.cl/producto/webpay#codigos-de-respuesta-de-autorizacion
const UseErrorLevelTwo = false

var (
	levelOneErrorMap = map[int]string{
		-1: "possible error in transaction data entry",
		-2: "possible error on account balance or card capacity",
		-3: "generic transaction error",
		-4: "rejected by the issuer",
		-5: "possible fraud",
	}
	levelTwoErrorMap = map[int]string{
		-1: "invalid card",
		-2: "connection error",
		-3: "exceeds maximum amount",
		-4: "invalid expiration date",
		-5: "authentication problem",
		-6: "general rejection",
		-7: "locked card",
		-8: "expired card",
		-9: "transaction not supported",
		-10: "transaction problem",
	}
)

type UnexpectedError struct {
	Err error
}

func (e *UnexpectedError) Error() string { return fmt.Sprintf("unexpected error: %s", e.Err.Error()) }

func (e *UnexpectedError) Unwrap() error { return e.Err }

type HttpProtocolError struct {
	StatusCode int
	Response *http.Response
}

func (e *HttpProtocolError) Error() string { return fmt.Sprintf("request failed with status %d", e.StatusCode) }

type TransactionError struct {
	Text string
	IsLevelTwo bool
	Code int
}

func (e *TransactionError) Error() string{
	return e.Text
}

var (
	ErrUnexpected = &UnexpectedError{errors.New("unexpected error")}
	ErrHttpProtocol = &HttpProtocolError{}
	ErrTransaction = &TransactionError{
		Text: "unknown error code received",
		IsLevelTwo: false,
	}
)

func mapError(code int) error {
	if code == 0 {
		return nil
	}

	if UseErrorLevelTwo {
		ErrTransaction.IsLevelTwo = true
		txt, ok := levelTwoErrorMap[code]
		if ok {
			ErrTransaction.Text = txt
		}
		return ErrTransaction
	}

	txt, ok := levelOneErrorMap[code]
	if ok {
		ErrTransaction.Text = txt
	}
	ErrTransaction.Text = txt
	return ErrTransaction
}
