package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/nikhil-github/api-cab-data/pkg/output"
)

type DBQueryer interface {
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	Select(dest interface{}, query string, args ...interface{}) error
	Rebind(query string) string
}

// Queryer provides database query operations.
type Queryer struct {
	db     DBQueryer
	logger *zap.Logger
}

// NewQueryer returns a new instance to query cab trip data.
func NewQueryer(db DBQueryer, logger *zap.Logger) *Queryer {
	return &Queryer{db: db, logger: logger}
}

// TripsByMedallionsOnPickUpDate get the count of trips for a cab by medallion and pick up date.
func (q *Queryer) TripsByMedallionsOnPickUpDate(ctx context.Context, medallion string, pickUpDate time.Time) (output.Result, error) {
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
	rows, err := q.db.QueryxContext(ctx, query, medallion, pickUpDate)
	if err != nil {
		q.logger.Error("sql error on query", zap.Error(err))
		return output.Result{}, errors.Wrap(err, "failed to query")
	}
	defer rows.Close()
	var count int
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			q.logger.Error("sql error on scan", zap.Error(err))
			return output.Result{}, errors.Wrap(err, "failed to query count")
		}
	}
	err = rows.Err()
	if err != nil {
		q.logger.Error("sql error on rows", zap.Error(err))
		return output.Result{}, err
	}
	return output.Result{Medallion: medallion, Trips: count}, nil
}

// Trips get the count of trips for a cab by medallion.
func (q *Queryer) TripsByMedallion(ctx context.Context, medallions []string) ([]output.Result, error) {
	var res []output.Result
	rawQuery := `
		SELECT
			medallion,
			count(medallion) AS trips
		FROM
			cab_trip_data
		WHERE	
			medallion IN (?)
		GROUP BY medallion;
	`
	query, args, err := sqlx.In(rawQuery, medallions)
	query = q.db.Rebind(query)
	err = q.db.Select(&res, query, args...)
	fmt.Println(err)
	if err != nil {
		q.logger.Error("sql error on query", zap.Error(err))
		return nil, errors.Wrap(err, "failed to query")
	}
	return res, nil
}
