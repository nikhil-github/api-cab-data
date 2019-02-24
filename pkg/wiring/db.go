package wiring

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// NewDatabase creates a new database
func NewDatabase(config DBConfig) (*sqlx.DB, error) {
	db, err := sql.Open("mysql", string(config.URL))
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(config.Connections.Max)
	db.SetMaxIdleConns(config.Connections.Idle)
	db.SetConnMaxLifetime(config.Connections.Lifetime)

	dbx := sqlx.NewDb(db, "mysql")
	return dbx, nil
}
