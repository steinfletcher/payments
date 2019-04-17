package api_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	m "github.com/petergtz/pegomock"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	acme "github.com/steinfletcher/payments"
	"github.com/steinfletcher/payments/api"
	"github.com/steinfletcher/payments/mocks"
	"github.com/steinfletcher/payments/test"
	"github.com/stretchr/testify/assert"
)

func TestCreatePayment_Success(t *testing.T) {
	id := uuid.New()
	var payment acme.Payment
	readJSON("testdata/create_payment.json", &payment)

	paymentService := mocks.NewMockPaymentService()
	m.When(paymentService.Create(payment)).ThenReturn(id, nil)

	apiTest(paymentService).
		Post("/v1/payment").
		JSON(readFile("testdata/create_payment.json")).
		Expect(t).
		Status(http.StatusCreated).
		Header("Location", id.String()).
		End()
}

func TestCreatePayment_InvalidAttributes(t *testing.T) {
	apiTest(mocks.NewMockPaymentService()).
		Post("/v1/payment").
		JSON(readFile("testdata/create_payment_with_invalid_attributes.json")).
		Expect(t).
		Status(http.StatusBadRequest).
		HeaderNotPresent("Location").
		Body(`{
			"code": "INVALID_FIELD",
			"detail": "invalid attributes: [(root): currency is required]"
		}`).
		End()
}

func TestCreatePayment_WithoutMandatoryField(t *testing.T) {
	apiTest(mocks.NewMockPaymentService()).
		Post("/v1/payment").
		JSON(readFile("testdata/create_payment_without_mandatory_field.json")).
		Expect(t).
		Status(http.StatusBadRequest).
		HeaderNotPresent("Location").
		Body(`{
			"code": "INVALID_FIELD",
			"detail": "organisation Id must be provided"
		}`).
		End()
}

func TestCreatePayment_InvalidRequestBody(t *testing.T) {
	apiTest(mocks.NewMockPaymentService()).
		Post("/v1/payment").
		JSON(readFile("testdata/invalid_request_body.json")).
		Expect(t).
		Status(http.StatusBadRequest).
		HeaderNotPresent("Location").
		Body(`{
			"code": "INVALID_REQUEST_BODY",
			"detail": "The request body is not valid"
		}`).
		End()
}

func TestCreatePayment_WithServiceError(t *testing.T) {
	var payment acme.Payment
	readJSON("testdata/create_payment.json", &payment)
	paymentService := mocks.NewMockPaymentService()
	m.When(paymentService.Create(payment)).ThenReturn(uuid.Nil, acme.ServerError)

	apiTest(paymentService).
		Post("/v1/payment").
		JSON(readFile("testdata/create_payment.json")).
		Expect(t).
		Status(http.StatusInternalServerError).
		HeaderNotPresent("Location").
		Body(`{
			"code": "SERVER_ERROR",
			"detail": "Sorry, something went wrong"
		}`).
		End()
}

func TestGetPayment_Success(t *testing.T) {
	id := uuid.New()
	paymentService := mocks.NewMockPaymentService()
	m.When(paymentService.Get(id)).ThenReturn(aPayment(id), nil)

	apiTest(paymentService).
		Get(fmt.Sprintf("/v1/payment/%s", id)).
		Expect(t).
		Status(http.StatusOK).
		Body(fmt.Sprintf(`{
			"id": "%s",
			"version": 0,
			"organisation_id": "57a3b643-cf4f-4f70-8636-0ddcdec07d68",
			"attributes": {
				"key": "value"
			}
		}`, id)).
		End()
}

func TestGetPayment_NotFound(t *testing.T) {
	id := uuid.New()
	paymentService := mocks.NewMockPaymentService()
	m.When(paymentService.Get(id)).ThenReturn(acme.Payment{}, acme.PaymentNotFound)

	apiTest(paymentService).
		Get(fmt.Sprintf("/v1/payment/%s", id)).
		Expect(t).
		Status(http.StatusBadRequest).
		Body(`{
			"code": "PAYMENT_NOT_FOUND",
			"detail": "We could not find a payment with the given ID"
		}`).
		End()
}

func TestGetPayment_InvalidID(t *testing.T) {
	apiTest(mocks.NewMockPaymentService()).
		Get("/v1/payment/invalidID").
		Expect(t).
		Status(http.StatusBadRequest).
		Body(`{
			"code": "INVALID_PAYMENT_ID",
			"detail": "The provided ID is not valid"
		}`).
		End()
}

func TestGetAllPayments_Success(t *testing.T) {
	id := uuid.New()
	paymentService := mocks.NewMockPaymentService()
	m.When(paymentService.GetAll()).ThenReturn(acme.Payments{
		Data: []acme.Payment{
			aPayment(id),
		},
	}, nil)

	apiTest(paymentService).
		Get("/v1/payment").
		Expect(t).
		Body(fmt.Sprintf(`{
			"data": [{
				"attributes": {
					"key": "value"
				},
				"id": "%s",
				"organisation_id": "57a3b643-cf4f-4f70-8636-0ddcdec07d68",
				"version": 0
			}]
		}`, id)).
		Status(http.StatusOK).
		End()
}

func TestGetAllPayments_EmptyArrayIfNone(t *testing.T) {
	paymentService := mocks.NewMockPaymentService()
	m.When(paymentService.GetAll()).ThenReturn(acme.Payments{
		Data: []acme.Payment{},
	}, nil)

	apiTest(paymentService).
		Get("/v1/payment").
		Expect(t).
		Body(fmt.Sprintf(`{
		  "data": []
		}`)).
		Status(http.StatusOK).
		End()
}

func TestGetAllPayments_ServerError(t *testing.T) {
	paymentService := mocks.NewMockPaymentService()
	m.When(paymentService.GetAll()).ThenReturn(acme.Payments{}, acme.ServerError)

	apiTest(paymentService).
		Get("/v1/payment").
		Expect(t).
		Body(`{
			"code": "SERVER_ERROR",
			"detail": "Sorry, something went wrong"
		}`).
		Status(http.StatusInternalServerError).
		End()
}

func TestDeletePayment_Success(t *testing.T) {
	id := uuid.New()
	paymentService := mocks.NewMockPaymentService()
	m.When(paymentService.Delete(id)).ThenReturn(nil)

	apiTest(paymentService).Debug().
		Delete(fmt.Sprintf("/v1/payment/%s", id)).
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestDeletePayment_WithInvalidID(t *testing.T) {
	apiTest(mocks.NewMockPaymentService()).
		Delete("/v1/payment/invalid_uuid").
		Expect(t).
		Status(http.StatusBadRequest).
		Assert(jsonpath.Equal("$.code", acme.InvalidID.Code)).
		End()
}

func TestDeletePayment_NotFound(t *testing.T) {
	id := uuid.New()
	paymentService := mocks.NewMockPaymentService()
	m.When(paymentService.Delete(id)).ThenReturn(acme.PaymentNotFound)

	apiTest(paymentService).
		Delete(fmt.Sprintf("/v1/payment/%s", id)).
		Expect(t).
		Status(http.StatusBadRequest).
		Assert(jsonpath.Equal("$.code", acme.PaymentNotFound.Code)).
		End()
}

func TestUpdatePayment_Success(t *testing.T) {
	id := uuid.New()
	var payment acme.Payment
	readJSON("testdata/update_payment.json", &payment)

	paymentService := mocks.NewMockPaymentService()
	m.When(paymentService.Update(id, payment)).ThenReturn(nil)

	apiTest(paymentService).
		Put(fmt.Sprintf("/v1/payment/%s", id)).
		JSON(readFile("testdata/update_payment.json")).
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestUpdatePayment_InvalidID(t *testing.T) {
	apiTest(mocks.NewMockPaymentService()).
		Put("/v1/payment/not_a_uuid").
		JSON(readFile("testdata/update_payment.json")).
		Expect(t).
		Status(http.StatusBadRequest).
		Body(`{
			"code": "INVALID_PAYMENT_ID",
			"detail": "The provided ID is not valid"
		}`).
		End()
}

func TestUpdatePayment_InvalidRequestBody(t *testing.T) {
	apiTest(mocks.NewMockPaymentService()).
		Put(fmt.Sprintf("/v1/payment/%s", uuid.New())).
		JSON(readFile("testdata/invalid_request_body.json")).
		Expect(t).
		Status(http.StatusBadRequest).
		Body(`{
			"code": "INVALID_FIELD",
			"detail": "The request body is not valid"
		}`).
		End()
}

func TestUpdatePayment_ServiceError(t *testing.T) {
	id := uuid.New()
	var payment acme.Payment
	readJSON("testdata/update_payment.json", &payment)

	paymentService := mocks.NewMockPaymentService()
	m.When(paymentService.Update(id, payment)).ThenReturn(acme.ServerError)

	apiTest(paymentService).
		Put(fmt.Sprintf("/v1/payment/%s", id)).
		JSON(readFile("testdata/update_payment.json")).
		Expect(t).
		Status(http.StatusInternalServerError).
		Body(`{
			"code": "SERVER_ERROR",
			"detail": "Sorry, something went wrong"
		}`).
		End()
}

func TestHealthCheck(t *testing.T) {
	srv := api.NewServer(mocks.NewMockPaymentService())
	go srv.Start("9001")
	defer srv.Close()
	cli := http.Client{Timeout: 1 * time.Second}

	res, err := cli.Get("http://localhost:9001/health")
	assert.NoError(t, err)
	data, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.JSONEq(t, `{"status":"OK"}`, string(data))
}

func aPayment(id ...uuid.UUID) acme.Payment {
	payment := acme.Payment{
		Version:        0,
		OrganisationID: uuid.MustParse("57a3b643-cf4f-4f70-8636-0ddcdec07d68"),
		Attributes:     map[string]string{"key": "value"},
	}
	if len(id) > 0 {
		payment.ID = id[0]
	}
	return payment
}

func apiTest(service acme.PaymentService) *apitest.APITest {
	return apitest.New().
		Recorder(test.Recorder).
		Handler(api.NewServer(service).Router)
}

func readFile(file string) string {
	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return string(fileContent)
}

func readJSON(file string, data interface{}) {
	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(fileContent, data)
	if err != nil {
		panic(err)
	}
}
