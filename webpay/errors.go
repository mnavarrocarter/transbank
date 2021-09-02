package webpay

// UseErrorLevelTwo indicates whether to map errors on level 1 or 2.
// See https://www.transbankdevelopers.cl/producto/webpay#codigos-de-respuesta-de-autorizacion
const UseErrorLevelTwo = false

type ErrTransaction int

func (e ErrTransaction) Error() string {
	if UseErrorLevelTwo {
		txt, ok := levelTwoErrorMap[int(e)]
		if ok {
			return txt
		}
	}

	txt, ok := levelOneErrorMap[int(e)]
	if ok {
		return txt
	}
	return "unknown transaction error"
}

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
