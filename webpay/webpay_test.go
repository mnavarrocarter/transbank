package webpay_test

// All test requests should have these headers
var reqHeaders = map[string]string{
	"Tbk-Api-Key-Id":     "597055555532",
	"Tbk-Api-Key-Secret": "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C",
	"Content-Type":       "application/json; charset=utf8",
	"Accept":             "application/json; charset=utf8",
}

// All test responses should have these headers
var resHeaders = map[string]string{
	"Content-Type":       "application/json",
}
