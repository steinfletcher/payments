package acme

import (
	"github.com/google/uuid"
)

//go:generate pegomock generate --use-experimental-model-gen --output-dir mocks PaymentService

type PaymentService interface {
	Get(id uuid.UUID) (Payment, error)
	GetAll() (Payments, error)
	Delete(id uuid.UUID) error
	Update(id uuid.UUID, payment Payment) error
	Create(payment Payment) (uuid.UUID, error)
}

type Payment struct {
	ID             uuid.UUID   `json:"id"`
	Version        int         `json:"version"`
	OrganisationID uuid.UUID   `json:"organisation_id"`
	Attributes     interface{} `json:"attributes"`
}

type Payments struct {
	Data []Payment `json:"data"`
}
