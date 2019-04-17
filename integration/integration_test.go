package integration_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/steinfletcher/apitest"
	"github.com/steinfletcher/payments/api"
	"github.com/steinfletcher/payments/postgres"
	"github.com/steinfletcher/payments/test"
)

func TestGetAllPayments(t *testing.T) {
	setup(t)

	paymentID := withPayment(t)

	apiTest().
		Get("/v1/payment").
		Expect(t).
		Body(fmt.Sprintf(readFile("testdata/expected_payments_response.json"), paymentID)).
		Status(http.StatusOK).
		End()
}

func TestGetPayment(t *testing.T) {
	setup(t)

	paymentID := withPayment(t)

	apiTest().
		Get(fmt.Sprintf("/v1/payment/%s", paymentID)).
		Expect(t).
		Body(fmt.Sprintf(readFile("testdata/expected_payment_response.json"), paymentID)).
		Status(http.StatusOK).
		End()
}

func TestCreatePayment(t *testing.T) {
	setup(t)

	apiTest().
		Post("/v1/payment").
		JSON(readFile("testdata/create_payment.json")).
		Expect(t).
		Status(http.StatusCreated).
		HeaderPresent("Location").
		End()
}

func TestDeletePayment(t *testing.T) {
	setup(t)

	paymentID := withPayment(t)

	apiTest().
		Delete(fmt.Sprintf("/v1/payment/%s", paymentID)).
		Expect(t).
		Status(http.StatusOK).
		End()

	apiTest().
		Get(fmt.Sprintf("/v1/payment/%s", paymentID)).
		Expect(t).
		Status(http.StatusBadRequest).
		Body(`{"code":"PAYMENT_NOT_FOUND", "detail":"We could not find a payment with the given ID"}`).
		End()
}

func withPayment(t *testing.T) string {
	result := apiTest().
		Post("/v1/payment").
		JSON(readFile("testdata/create_payment.json")).
		Expect(t).
		Status(http.StatusCreated).
		HeaderPresent("Location").
		End()
	return result.Response.Header.Get("Location")
}

func setup(t *testing.T) {
	if os.Getenv("DB_ADDR") == "" {
		t.Skip()
	}
	// teardown
	test.DBSetup(func(tx *sqlx.Tx) {})
}

func apiTest() *apitest.APITest {
	service := postgres.NewPaymentRepository(test.DBConnect())
	return apitest.New().
		Recorder(test.Recorder).
		Handler(api.NewServer(service).Router).
		Report(apitest.SequenceDiagram())
}

func readFile(file string) string {
	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return string(fileContent)
}
