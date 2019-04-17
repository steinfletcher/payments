# payments

## Prerequisites

- docker
- Go 1.12

## Design

The API exposes the following endpoints.

* `GET    /v1/payment`      All payments
* `GET    /v1/payment/:id`  Individual payment by ID
* `POST   /v1/payment`      Create payment
* `PUT    /v1/payment/:id`  Update payment by ID
* `DELETE /v1/payment/:id`  Delete payment by ID

### Package layout

The package layout strategy is based on 3 simple rules:

* Root package is for domain types
* Group subpackages by dependency
* Main package ties together dependencies (IoC)

These rules help decouple concerns and eliminate circular dependencies. It helps to define a clear domain language across the application.

See this [article](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1) for more information behind this approach.

### Persistence

Postgres is used for persistence due to its rich feature set - transaction support, structured and unstructured (JSON) types, support by major cloud vendors, e.g. RDS in AWS. Local development workflow is smooth since Postgres can be run using docker.

The following decisions were made when designing the persistence layer

* Immutable data store. `SELECT` and `INSERT` are the only database operations supported. Deletes at the API level are soft deletes in the database - a new deletion record is created. This generates an audit trail. The same approach is taken for updates - a new database row is inserted and a version number is incremented. The API read operations retrieve the latest version.

* Attributes stored as JSON. JSON schema is used to validate the attributes before inserting into the database ensuring data integrity. These attributes might be modelled as first class entities in future as requirements are defined. Storing as JSON is a quick and simple approach suitable for an initial version.

## Tests

### Unit

`make test`

The definition of a `unit` in the project is a package. Units are tested using the `<package>_test` convention to ensure tests run against the public API of the package and not implementation details. This helps prevent brittle tests that fail when implementation details change.

Unit tests generally mock external collaborators and dependency injection is used throughout to achieve this.

### Integration

```bash
make postgres-dev
make test-integration # OR make test-all to include unit tests
```

These tests ensure that the units are working together. A real postgres is used to ensure the database interactions are correct. The `postgres` package is tested in isolation. 

The `integration` package contains behavioural tests. These tests send a HTTP request to the application and assert on an expected response. These tests also render sequence diagrams on completion.

## Run locally

See the makefile for commands to run unit and integration tests and the build. 
For local development the api can be run with postgres in docker compose as follows

```bash
docker-compose up
```

Create a payment using curl

```bash
curl -X POST \
  http://localhost:9000/v1/payment \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
  "organisation_id": "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb",
  "attributes": {
    "amount": "100.21",
    "beneficiary_party": {
      "account_name": "W Owens",
      "account_number": "31926819",
      "account_number_code": "BBAN",
      "account_type": 0,
      "address": "1 The Beneficiary Localtown SE2",
      "bank_id": "403000",
      "bank_id_code": "GBDSC",
      "name": "Wilfred Jeremiah Owens"
    },
    "charges_information": {
      "bearer_code": "SHAR",
      "sender_charges": [
        {
          "amount": "5.00",
          "currency": "GBP"
        },
        {
          "amount": "10.00",
          "currency": "USD"
        }
      ],
      "receiver_charges_amount": "1.00",
      "receiver_charges_currency": "USD"
    },
    "currency": "GBP",
    "debtor_party": {
      "account_name": "EJ Brown Black",
      "account_number": "GB29XABC10161234567801",
      "account_number_code": "IBAN",
      "address": "10 Debtor Crescent Sourcetown NE1",
      "bank_id": "203301",
      "bank_id_code": "GBDSC",
      "name": "Emelia Jane Brown"
    },
    "end_to_end_reference": "Wil piano Jan",
    "fx": {
      "contract_reference": "FX123",
      "exchange_rate": "2.00000",
      "original_amount": "200.42",
      "original_currency": "USD"
    },
    "numeric_reference": "1002001",
    "payment_id": "123456789012345678",
    "payment_purpose": "Paying for goods/services",
    "payment_scheme": "FPS",
    "payment_type": "Credit",
    "processing_date": "2017-01-18",
    "reference": "Payment for Em'\''s piano lessons",
    "scheme_payment_sub_type": "InternetBanking",
    "scheme_payment_type": "ImmediatePayment",
    "sponsor_party": {
      "account_number": "56781234",
      "bank_id": "123123",
      "bank_id_code": "GBDSC"
    }
  }
}
'
```

Then read the payment 

```bash
curl -X GET \
  http://localhost:9000/v1/payment \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache'
```
