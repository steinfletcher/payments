package postgres

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin/json"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/pkg/errors"
	"github.com/steinfletcher/payments"
)

const getQuery = `SELECT p.version, p.external_id, p.organisation_id, p.attributes
FROM payments p
         JOIN (
    SELECT MAX(version) as version, external_id
    FROM payments vp
    GROUP BY external_id) t
              ON t.external_id = p.external_id AND t.version = p.version
WHERE p.deleted = FALSE %s`

const insertQuery = `INSERT INTO payments (external_id, attributes, organisation_id, version, deleted)
 VALUES ($1, $2, $3, $4, $5)`

type paymentRepository struct {
	db *sqlx.DB
}

type paymentRecord struct {
	Version        int            `db:"version"`
	ExternalID     string         `db:"external_id"`
	OrganisationID string         `db:"organisation_id"`
	Attributes     types.JSONText `db:"attributes"`
}

func (r *paymentRepository) GetAll() (acme.Payments, error) {
	var p []paymentRecord
	err := withTx(r.db, func(tx *sql.Tx) error {
		return r.db.Select(&p, fmt.Sprintf(getQuery, ""))
	})
	if err != nil {
		return acme.Payments{}, err
	}
	return mapPayments(p), nil
}

func (r *paymentRepository) Get(id uuid.UUID) (acme.Payment, error) {
	var p paymentRecord
	err := withTx(r.db, func(tx *sql.Tx) error {
		err := r.db.Get(&p, fmt.Sprintf(getQuery, fmt.Sprintf("AND p.external_id = '%s'", id)))
		if err != nil {
			if err == sql.ErrNoRows {
				return acme.PaymentNotFound
			}
			return errors.WithStack(acme.ServerError)
		}
		return nil
	})
	if err != nil {
		return acme.Payment{}, err
	}
	return mapPayment(p), nil
}

func (r *paymentRepository) Create(p acme.Payment) (uuid.UUID, error) {
	newID := uuid.New()
	attributes, err := json.Marshal(p.Attributes)
	if err != nil {
		return newID, acme.ServerError
	}

	err = withTx(r.db, func(tx *sql.Tx) error {
		_, err := r.db.Exec(insertQuery, newID, attributes, p.OrganisationID, p.Version, false)
		if err != nil {
			return errors.WithStack(acme.ServerError)
		}
		return nil
	})
	return newID, err
}

func (r *paymentRepository) Update(id uuid.UUID, updatedPayment acme.Payment) error {
	return withTx(r.db, func(tx *sql.Tx) error {
		payment, err := r.Get(id)
		if err != nil {
			return err
		}

		result, err := r.db.Exec(insertQuery, payment.ID, updatedPayment.Attributes, updatedPayment.OrganisationID,
			payment.Version+1, false)
		if err != nil {
			return errors.WithStack(acme.ServerError)
		}

		rows, _ := result.RowsAffected()
		if rows <= 0 {
			return acme.PaymentNotFound
		}
		return nil
	})
}

func (r *paymentRepository) Delete(id uuid.UUID) error {
	return withTx(r.db, func(tx *sql.Tx) error {
		payment, err := r.Get(id)
		if err != nil {
			return err
		}

		result, err := r.db.Exec(insertQuery, payment.ID, payment.Attributes, payment.OrganisationID, payment.Version+1, true)
		if err != nil {
			return errors.WithStack(acme.ServerError)
		}

		rows, _ := result.RowsAffected()
		if rows <= 0 {
			return acme.PaymentNotFound
		}
		return nil
	})
}

func NewPaymentRepository(db *sqlx.DB) acme.PaymentService {
	return &paymentRepository{db}
}

func mapPayments(dbRecords []paymentRecord) acme.Payments {
	payments := []acme.Payment{}
	for _, v := range dbRecords {
		payments = append(payments, mapPayment(v))
	}
	return acme.Payments{
		Data: payments,
	}
}

func mapPayment(dbRecord paymentRecord) acme.Payment {
	return acme.Payment{
		ID:             uuid.MustParse(dbRecord.ExternalID),
		Version:        dbRecord.Version,
		OrganisationID: uuid.MustParse(dbRecord.OrganisationID),
		Attributes:     dbRecord.Attributes,
	}
}

// withTx encapsulates transaction concerns such as rollbacks and commit.
// This helps decouple lower level transaction handling from business logic.
func withTx(db *sqlx.DB, fn func(*sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return errors.WithStack(acme.ServerError)
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
