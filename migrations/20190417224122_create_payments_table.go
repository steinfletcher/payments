package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up20190417224122, Down20190417224122)
}

func Up20190417224122(tx *sql.Tx) error {
	return exec(`CREATE TABLE payments
(
    id              SERIAL PRIMARY KEY NOT NULL,
    version         INT                NOT NULL DEFAULT 0,
    external_id     TEXT               NOT NULL,
    attributes      JSON               NOT NULL,
    organisation_id TEXT               NOT NULL,
    deleted         BOOLEAN                     DEFAULT FALSE
);

CREATE UNIQUE INDEX payments_external_id ON payments (external_id, version);
`, tx)
}

func Down20190417224122(tx *sql.Tx) error {
	return exec("DROP TABLE payments;", tx)
}

func exec(query string, tx *sql.Tx) error {
	_, err := tx.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
