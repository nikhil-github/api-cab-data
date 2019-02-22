package wiring

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

// NewDatabase creates a new database
func NewDatabase(config DatabaseConfig) (*sqlx.DB, error) {
	// db, err := sqlx.Open("mysql", string("root:password@tcp(localhost:3306)/cabtrips?parseTime=true"))
	db, err := sql.Open("mysql", string(config.URL))
	if err != nil {
		return nil,  err
	}
	db.SetMaxOpenConns(config.Connections.Max)
	db.SetMaxIdleConns(config.Connections.Idle)
	db.SetConnMaxLifetime(config.Connections.Lifetime)

	dbx := sqlx.NewDb(db, "postgres")
	return dbx, nil
}
