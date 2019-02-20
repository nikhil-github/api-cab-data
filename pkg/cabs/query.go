package cabs

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type DBQueryer interface {
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
}

// Queryer is for providing database query operations.
type Queryer struct {
	db   DBQueryer
}

// NewQueryer returns a new instance to query cab trip data.
func NewQueryer(db DBQueryer) *Queryer {
	return &Queryer{db: db}
}

// CabTripsByPickUpDate get the count of trips for a cab by pick up date.
func (q *Queryer) CabTripsByPickUpDate(ctx context.Context, medallion string,pickUpDate time.Time) (int, error) {
	query := `
		SELECT
			count(*)
		AS
			count
		FROM
			cab_trip_data
		WHERE	
			medallion = $1
		AND
			pickup_datetime::DATE = $2::DATE;
	`
	rows,err := q.db.QueryxContext(ctx, query, medallion,pickUpDate)
	if err != nil {
		return 0,errors.Wrap(err,"failed to query")
	}
	defer rows.Close()
	var count int
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0,errors.Wrap(err,"failed to query count")
		}
	}
	err = rows.Err()
	if err != nil {
		return 0, err
	}
	return count, nil
}
