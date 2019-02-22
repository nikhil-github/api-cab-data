package database

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type DBQueryer interface {
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
}

// Queryer is for providing database query operations.
type Queryer struct {
	db DBQueryer
	logger *zap.Logger
}

// NewQueryer returns a new instance to query cab trip data.
func NewQueryer(db DBQueryer,logger *zap.Logger) *Queryer {
	return &Queryer{db: db,logger:logger}
}

// TripsByPickUpDate get the count of trips for a cab by pick up date.
func (q *Queryer) TripsByPickUpDate(ctx context.Context, medallion string,pickUpDate time.Time) (int, error) {
	query := `
		SELECT
			count(medallion)
		AS
			count
		FROM
			cab_trip_data
		WHERE	
			medallion = ?
		AND
			DATE(pickup_datetime) = DATE(?)
	`
	rows,err := q.db.QueryxContext(ctx, query, medallion,pickUpDate)
	if err != nil {
		q.logger.Error("sql error on query",zap.Error(err))
		return 0,errors.Wrap(err,"failed to query")
	}
	defer rows.Close()
	var count int
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			q.logger.Error("sql error on scan",zap.Error(err))
			return 0,errors.Wrap(err,"failed to query count")
		}
	}
	err = rows.Err()
	if err != nil {
		q.logger.Error("sql error on rows",zap.Error(err))
		return 0, err
	}
	return count, nil
}
