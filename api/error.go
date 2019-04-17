package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	acme "github.com/steinfletcher/payments"
)

// errorToStatusCodeLookup maps application errors to http status codes
// It helps to decouple application errors from the HTTP layer
var errorToStatusCodeLookup = map[string]int{
	acme.InvalidID.Code:          http.StatusBadRequest,
	acme.InvalidRequestBody.Code: http.StatusBadRequest,
	acme.PaymentNotFound.Code:    http.StatusBadRequest,
	acme.InvalidField.Code:       http.StatusBadRequest,
	acme.ServerError.Code:        http.StatusInternalServerError,
}

// errorHandler is a middleware that sets any present application errors on the response
func errorHandler(c *gin.Context) {
	c.Next()

	if len(c.Errors) == 0 {
		return
	}

	for _, err := range c.Errors {
		if parsedError, ok := err.Err.(acme.Error); ok {
			statusCode := errorToStatusCodeLookup[parsedError.Code]
			c.JSON(statusCode, parsedError)
			return
		}
	}

	c.JSON(http.StatusInternalServerError, acme.ServerError)
}
