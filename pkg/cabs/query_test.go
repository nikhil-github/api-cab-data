package cabs_test

import (
	_ "github.com/lib/pq"
	"testing"
	"fmt"
	"github.com/jmoiron/sqlx"
	"bitbucket.org/ffxblue/api-video/lib/image"
	"time"
)

func TestCountByPickUpDate(t *testing.T) {
	type args struct {
		Limit   int
		SinceID int
	}
	type fields struct {

	}
	type want struct {
		Err     error
		VideoIDs  []int
		Message string
	}
	testTable := []struct {
		Name   string
		Args   args
		Fields fields
		Want   want
	}{
		{
			Name: "Success - get one",
			Args: args{Limit: 1},
		},
	}
	for _, tt := range testTable {
		t.Run(tt.Name, func(t *testing.T) {
			db, err := sqlx.Open("postgres", string("postgres://app:password@localhost:32768/app?sslmode=disable"))
			if err != nil {
				fmt.Println("failed DB")
			}
			d,_ := time.ParseDuration("100h")
			qs := image.NewQueryer(db,d, tracer)

			ids,err := qs.Unprocessed(ctx,0,3)
			fmt.Println("err=>",err)
			fmt.Println("records=>",ids)

		})
	}

}
