package acme

var InvalidID = Error{
	Code:   "INVALID_PAYMENT_ID",
	Detail: "The provided ID is not valid",
}

var ServerError = Error{
	Code:   "SERVER_ERROR",
	Detail: "Sorry, something went wrong",
}

var PaymentNotFound = Error{
	Code:   "PAYMENT_NOT_FOUND",
	Detail: "We could not find a payment with the given ID",
}

var InvalidField = Error{
	Code:   "INVALID_FIELD",
	Detail: "The request body is not valid",
}

var InvalidRequestBody = Error{
	Code:   "INVALID_REQUEST_BODY",
	Detail: "The request body is not valid",
}

type Error struct {
	Code   string      `json:"code"`
	Detail string      `json:"detail"`
	Meta   interface{} `json:"meta,omitempty"`
}

func (r Error) Error() string {
	return r.Code
}
