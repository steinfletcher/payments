package postgres_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/steinfletcher/payments"
	_ "github.com/steinfletcher/payments/migrations"
	"github.com/steinfletcher/payments/postgres"
	"github.com/steinfletcher/payments/test"
	"github.com/stretchr/testify/assert"
)

func TestCreatePayment(t *testing.T) {
	test.SkipIntegration(t)
	db := test.DBSetup(func(tx *sqlx.Tx) {})
	organisationID := uuid.New()
	repository := postgres.NewPaymentRepository(db)

	id, err := repository.Create(acme.Payment{
		OrganisationID: organisationID,
		Attributes:     types.JSONText(`{"key":"value"}`),
		Version:        0,
	})
	assert.NoError(t, err)
	payment, err := repository.Get(id)

	assert.NoError(t, err)
	assert.Equal(t, acme.Payment{
		ID:             id,
		Version:        0,
		OrganisationID: organisationID,
		Attributes:     types.JSONText(`{"key":"value"}`),
	}, payment)
}

func TestUpdatePayment(t *testing.T) {
	test.SkipIntegration(t)
	externalID := uuid.New()
	organisationID := uuid.New()
	updatedOrganisationID := uuid.MustParse("efa9c7a2-5bc7-461b-a977-20853c7221cd")
	db := test.DBSetup(func(tx *sqlx.Tx) {
		query := `INSERT INTO payments (external_id, attributes, version, organisation_id) VALUES
				('%s', '{"key": "value"}', 0, '%s')`
		tx.MustExec(fmt.Sprintf(query, externalID, organisationID))
	})
	repository := postgres.NewPaymentRepository(db)

	err := repository.Update(externalID, acme.Payment{
		OrganisationID: updatedOrganisationID,
		Attributes:     types.JSONText(`{"key":"newValue"}`),
	})
	assert.NoError(t, err)
	payment, err := repository.Get(externalID)

	assert.NoError(t, err)
	assert.Equal(t, acme.Payment{
		ID:             externalID,
		Version:        1,
		OrganisationID: updatedOrganisationID,
		Attributes:     types.JSONText(`{"key":"newValue"}`),
	}, payment)
}

func TestDeletePayment_MarksThePaymentDeleted(t *testing.T) {
	test.SkipIntegration(t)
	externalID := uuid.New()
	db := test.DBSetup(func(tx *sqlx.Tx) {
		query := `INSERT INTO payments (external_id, attributes, version, organisation_id) VALUES
				('%s', '{"key": "value"}', 0, '%s')`
		tx.MustExec(fmt.Sprintf(query, externalID, uuid.New()))
	})
	payments := postgres.NewPaymentRepository(db)

	err := payments.Delete(externalID)

	assert.NoError(t, err)
	_, err = payments.Get(externalID)
	assert.EqualError(t, err, acme.PaymentNotFound.Code)
}

func TestDeletePayment_ReportsPaymentNotFound(t *testing.T) {
	test.SkipIntegration(t)
	db := test.DBSetup(func(tx *sqlx.Tx) {})

	err := postgres.NewPaymentRepository(db).Delete(uuid.New())

	assert.EqualError(t, err, acme.PaymentNotFound.Code)
}

func TestGetPayments_EmptyListIfNoPayments(t *testing.T) {
	test.SkipIntegration(t)
	db := test.DBSetup(func(tx *sqlx.Tx) {})

	payments, err := postgres.NewPaymentRepository(db).GetAll()

	assert.NoError(t, err)
	assert.Empty(t, payments.Data)
}

func TestGetPayments_DoesNotReturnDeletedPayments(t *testing.T) {
	test.SkipIntegration(t)
	db := test.DBSetup(func(tx *sqlx.Tx) {
		query := `INSERT INTO payments (external_id, attributes, version, organisation_id, deleted) VALUES
				('%s', '{"key": "valueOriginal"}', 0, '%s', true)`
		tx.MustExec(fmt.Sprintf(query, uuid.New(), uuid.New()))
	})

	payments, err := postgres.NewPaymentRepository(db).GetAll()

	assert.NoError(t, err)
	assert.Empty(t, payments.Data)
}

func TestGetPayments_ReturnsLatestVersion(t *testing.T) {
	test.SkipIntegration(t)
	externalID := uuid.New()
	organisationID := uuid.New()
	db := test.DBSetup(func(tx *sqlx.Tx) {
		v0 := `INSERT INTO payments (external_id, attributes, version, organisation_id) VALUES
				('%s', '{"key":"valueOriginal"}', 0, '%s')`
		tx.MustExec(fmt.Sprintf(v0, externalID, organisationID))

		v1 := `INSERT INTO payments (external_id, attributes, version, organisation_id) VALUES
				('%s', '{"key":"valueUpdated"}', 1, '%s')`
		tx.MustExec(fmt.Sprintf(v1, externalID, organisationID))
	})

	payments, err := postgres.NewPaymentRepository(db).GetAll()

	assert.NoError(t, err)
	assert.Equal(t, acme.Payments{
		Data: []acme.Payment{
			{
				ID:             externalID,
				Version:        1,
				OrganisationID: organisationID,
				Attributes:     types.JSONText(`{"key":"valueUpdated"}`),
			},
		},
	}, payments)
}

func TestGetPayment_ByID(t *testing.T) {
	test.SkipIntegration(t)
	externalID := uuid.New()
	organisationID := uuid.New()
	db := test.DBSetup(func(tx *sqlx.Tx) {
		query := `INSERT INTO payments (external_id, attributes, version, organisation_id) VALUES
				('%s', '{"key":"value"}', 0, '%s')`
		tx.MustExec(fmt.Sprintf(query, externalID, organisationID))
	})

	payment, err := postgres.NewPaymentRepository(db).Get(externalID)

	assert.NoError(t, err)
	assert.Equal(t, payment, acme.Payment{
		ID:             externalID,
		Version:        0,
		OrganisationID: organisationID,
		Attributes:     types.JSONText(`{"key":"value"}`),
	})
}

func TestGetPayment_ByID_NotFound(t *testing.T) {
	test.SkipIntegration(t)
	db := test.DBSetup(func(tx *sqlx.Tx) {})

	_, err := postgres.NewPaymentRepository(db).Get(uuid.New())

	assert.EqualError(t, err, acme.PaymentNotFound.Code)
}
