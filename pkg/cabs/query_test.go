package cabs_test

import (
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"testing"
	"time"
	"github.com/jmoiron/sqlx"
	"github.com/nikhil-github/api-cab-data/pkg/cabs"
	"context"
)

func TestCountByPickUpDate(t *testing.T) {
	// 2013-12-01 00:13:00
	date := time.Date(2013,12,31,0,1,0,0,time.UTC)
	type args struct {
		Medallion string
		PickUpDate  time.Time
	}
	type fields struct {

	}
	type want struct {
		Count int
		Err     error
	}
	testTable := []struct {
		Name   string
		Args   args
		Fields fields
		Want   want
	}{
		{
			Name: "Success - get one",
			Args: args{Medallion: "temp", PickUpDate: time.Time{}},
		},
	}
	for _, tt := range testTable {
		t.Run(tt.Name, func(t *testing.T) {
			db, err := sqlx.Open("mysql", string("root:password@tcp(localhost:3306)/cabtrips?parseTime=true"))
			if err != nil {
				fmt.Println("failed DB",err)
			}
			fmt.Println("database connetion ",db.DriverName())
			svc := cabs.NewQueryer(db)
			count, err := svc.CabTripsByPickUpDate(context.Background(),"67EB082BFFE72095EAF18488BEA96050",date)
			fmt.Println("err=>",err)
			fmt.Println("records=>",count)

		})
	}

}
