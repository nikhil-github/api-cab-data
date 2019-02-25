package database_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/nikhil-github/api-cab-data/pkg/database"
	"github.com/nikhil-github/api-cab-data/pkg/output"
)

func TestTripsByMedallionOnPickUpDate(t *testing.T) {
	pDate := time.Date(2013, 12, 31, 0, 1, 0, 0, time.UTC)
	type args struct {
		Medallion  string
		PickUpDate time.Time
	}
	type fields struct {
		MockOperations func(sqlmock.Sqlmock)
	}
	type want struct {
		Error  string
		Result output.Result
	}

	testTable := []struct {
		Name   string
		Args   args
		Fields fields
		Want   want
	}{
		{
			Name: "Success, record found",
			Args: args{Medallion: "67EB082BFFE72095EAF18488BEA96050", PickUpDate: pDate},
			Fields: fields{MockOperations: func(m sqlmock.Sqlmock) {
				columns := []string{"count"}
				rows := sqlmock.NewRows(columns)
				rows.AddRow(1)
				selectCount(m, "67EB082BFFE72095EAF18488BEA96050", pDate).WillReturnRows(rows)
			}},
			Want: want{Result: output.Result{Medallion: "67EB082BFFE72095EAF18488BEA96050", Trips: 1}},
		},
		{
			Name: "Failure, DB error",
			Args: args{Medallion: "55EB082BFFE795EAF18488BEA96050", PickUpDate: pDate},
			Fields: fields{MockOperations: func(m sqlmock.Sqlmock) {
				selectCount(m, "55EB082BFFE795EAF18488BEA96050", pDate).WillReturnError(errors.New("sql error"))
			}},
			Want: want{Error: "failed to query: sql error"},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.Name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err, "Unable to create Sqlmock DB")
			db := sqlx.NewDb(mockDB, "mysql")
			defer db.Close()
			tt.Fields.MockOperations(mock)

			dao := database.NewQueryer(db, zap.NewNop())
			res, err := dao.TripsByMedallionsOnPickUpDate(context.Background(), tt.Args.Medallion, tt.Args.PickUpDate)
			assert.NoError(t, mock.ExpectationsWereMet(), "DB Expectations")
			if tt.Want.Error != "" {
				assert.EqualError(t, err, tt.Want.Error, "Error")
				return
			}
			require.NoError(t, err, "Unexpected error")
			assert.Equal(t, tt.Want.Result, res, "Result")
		})
	}
}

func selectCount(m sqlmock.Sqlmock, medallion string, pDate time.Time) *sqlmock.ExpectedQuery {
	return m.ExpectQuery(`
		SELECT
			count\(medallion\)
		AS
			count
		FROM
			cab_trip_data
		WHERE	
			medallion = \?
		AND
			DATE\(pickup_datetime\) = DATE\(\?\)
	`).WithArgs(medallion, pDate)
}

func TestTripsByMedallion(t *testing.T) {
	type args struct {
		Medallions []string
	}
	type fields struct {
		MockOperations func(sqlmock.Sqlmock)
	}
	type want struct {
		Error  string
		Result []output.Result
	}

	testTable := []struct {
		Name   string
		Args   args
		Fields fields
		Want   want
	}{
		{
			Name: "Success, record found",
			Args: args{Medallions: []string{"67EB082BFFE72095EAF18488BEA96050"}},
			Fields: fields{MockOperations: func(m sqlmock.Sqlmock) {
				columns := []string{"medallion", "trips"}
				rows := sqlmock.NewRows(columns)
				rows.AddRow("67EB082BFFE72095EAF18488BEA96050", 1)
				selectCounts(m).WillReturnRows(rows)
			}},
			Want: want{Result: []output.Result{{Medallion: "67EB082BFFE72095EAF18488BEA96050", Trips: 1}}},
		},
		{
			Name: "Failure, DB error",
			Args: args{Medallions: []string{"55EB082BFFE795EAF18488BEA96050"}},
			Fields: fields{MockOperations: func(m sqlmock.Sqlmock) {
				selectCounts(m).WillReturnError(errors.New("sql error"))
			}},
			Want: want{Error: "failed to query: sql error"},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.Name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err, "Unable to create Sqlmock DB")
			db := sqlx.NewDb(mockDB, "mysql")
			defer db.Close()
			tt.Fields.MockOperations(mock)

			dao := database.NewQueryer(db, zap.NewNop())
			res, err := dao.TripsByMedallion(context.Background(), tt.Args.Medallions)
			assert.NoError(t, mock.ExpectationsWereMet(), "DB Expectations")
			if tt.Want.Error != "" {
				assert.EqualError(t, err, tt.Want.Error, "Error")
				return
			}
			require.NoError(t, err, "Unexpected error")
			assert.Equal(t, tt.Want.Result, res, "Result")
		})
	}
}

func selectCounts(m sqlmock.Sqlmock) *sqlmock.ExpectedQuery {
	return m.ExpectQuery(`
		SELECT
			medallion,
			count\(medallion\) AS trips
		FROM
			cab_trip_data
		WHERE	
			medallion IN \(\?\)
		GROUP BY medallion
	`)
}
