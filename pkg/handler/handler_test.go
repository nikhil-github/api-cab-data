package handler_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/nikhil-github/api-cab-data/pkg/output"
	"github.com/nikhil-github/api-cab-data/pkg/wiring"
)

func TestHandler_TripsByMedAndPickUpDate(t *testing.T) {
	res := output.Result{Medallion: "YYYY", Trips: 10}
	type args struct {
		URL  string
		Path string
	}
	type fields struct {
		MockExpectations func(m *mockTripSvc)
	}
	type want struct {
		Status int
		Body   string
	}
	testTable := []struct {
		Name   string
		Args   args
		Fields fields
		Want   want
	}{
		{
			Name:   "Failure - Invalid pickup date",
			Args:   args{Path: "/trips/v1/medallion/67EB082BFFE72095EAF18488BEA96050/pickupdate/201p-12-31"},
			Fields: fields{MockExpectations: func(m *mockTripSvc) {}},
			Want:   want{Status: http.StatusBadRequest},
		},
		{
			Name: "Service failed to query count",
			Args: args{Path: "/trips/v1/medallion/TTTTTTTT/pickupdate/2013-12-31"},
			Fields: fields{MockExpectations: func(m *mockTripSvc) {
				pd, err := time.Parse("2006-01-02", "2013-12-31")
				if err != nil {
					t.Fail()
				}
				m.OnTripsByPickUpDate("TTTTTTTT", pd, false).Return(output.Result{}, errors.New("error"))
			}},
			Want: want{Status: http.StatusInternalServerError},
		},
		{
			Name: "Success with one medallion",
			Args: args{Path: "/trips/v1/medallion/YYYY/pickupdate/2013-12-31"},
			Fields: fields{MockExpectations: func(m *mockTripSvc) {
				pd, err := time.Parse("2006-01-02", "2013-12-31")
				if err != nil {
					t.Fail()
				}
				m.OnTripsByPickUpDate("YYYY", pd, false).Return(res, nil)
			}},
			Want: want{Status: http.StatusOK, Body: `{"medallion":"YYYY","trips":10}`},
		},
		{
			Name: "Success with by pass cache flag",
			Args: args{Path: "/trips/v1/medallion/YYYY/pickupdate/2013-12-31?bypasscache=true"},
			Fields: fields{MockExpectations: func(m *mockTripSvc) {
				pd, err := time.Parse("2006-01-02", "2013-12-31")
				if err != nil {
					t.Fail()
				}
				m.OnTripsByPickUpDate("YYYY", pd, true).Return(res, nil)
			}},
			Want: want{Status: http.StatusOK, Body: `{"medallion":"YYYY","trips":10}`},
		},
	}
	for _, tt := range testTable {
		t.Run(tt.Name, func(t *testing.T) {
			logger := zap.NewNop()
			var m mockTripSvc
			tt.Fields.MockExpectations(&m)
			params := new(wiring.Params)
			params.Svc = &m
			params.Logger = logger
			mx := wiring.NewRouter(params)
			ts := httptest.NewServer(mx)
			defer ts.Close()
			res, err := http.Get(ts.URL + tt.Args.Path)
			assert.NoError(t, err, "Error executing request")
			defer res.Body.Close()
			m.AssertExpectations(t)
			assert.Equal(t, tt.Want.Status, res.StatusCode, "status")
			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err, "Error reading response")
			if tt.Want.Body != "" {
				assert.JSONEq(t, tt.Want.Body, string(body), "response")
			}
		})
	}
}

func TestHandler_TripsByMedallions(t *testing.T) {
	res := output.Result{Medallion: "YYYY", Trips: 10}
	res2 := output.Result{Medallion: "ZZZZ", Trips: 3}
	type args struct {
		URL  string
		Path string
	}
	type fields struct {
		MockExpectations func(m *mockTripSvc)
	}
	type want struct {
		Status int
		Body   string
	}
	testTable := []struct {
		Name   string
		Args   args
		Fields fields
		Want   want
	}{
		{
			Name: "Service failed to query count",
			Args: args{Path: "/trips/v1/medallions/TTTTTTTT"},
			Fields: fields{MockExpectations: func(m *mockTripSvc) {
				m.OnTripsTripsByMedallion([]string{"TTTTTTTT"}, false).Return([]output.Result{}, errors.New("error"))
			}},
			Want: want{Status: http.StatusInternalServerError},
		},
		{
			Name: "Success with one medallion",
			Args: args{Path: "/trips/v1/medallions/YYYY"},
			Fields: fields{MockExpectations: func(m *mockTripSvc) {

				m.OnTripsTripsByMedallion([]string{"YYYY"}, false).Return([]output.Result{res}, nil)
			}},
			Want: want{Status: http.StatusOK, Body: `[{"medallion":"YYYY","trips":10}]`},
		},
		{
			Name: "Success with multiple medallion",
			Args: args{Path: "/trips/v1/medallions/YYYY,ZZZZ"},
			Fields: fields{MockExpectations: func(m *mockTripSvc) {
				m.OnTripsTripsByMedallion([]string{"YYYY", "ZZZZ"}, false).Return([]output.Result{res, res2}, nil)
			}},
			Want: want{Status: http.StatusOK, Body: `[{"medallion":"YYYY","trips":10},{"medallion":"ZZZZ","trips":3}]`},
		},
	}
	for _, tt := range testTable {
		t.Run(tt.Name, func(t *testing.T) {
			logger := zap.NewNop()
			var m mockTripSvc
			tt.Fields.MockExpectations(&m)
			params := new(wiring.Params)
			params.Svc = &m
			params.Logger = logger
			mx := wiring.NewRouter(params)
			ts := httptest.NewServer(mx)
			defer ts.Close()
			res, err := http.Get(ts.URL + tt.Args.Path)
			assert.NoError(t, err, "Error executing request")
			defer res.Body.Close()
			m.AssertExpectations(t)
			assert.Equal(t, tt.Want.Status, res.StatusCode, "status")
			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err, "Error reading response")
			if tt.Want.Body != "" {
				assert.JSONEq(t, tt.Want.Body, string(body), "response")
			}
		})
	}
}

type mockTripSvc struct {
	mock.Mock
}

func (m *mockTripSvc) TripsByMedallionsOnPickUpDate(ctx context.Context, medallion string, pickUpDate time.Time, byPassCache bool) (output.Result, error) {
	args := m.Called(ctx, medallion, pickUpDate, byPassCache)
	return args.Get(0).(output.Result), args.Error(1)
}

func (m *mockTripSvc) OnTripsByPickUpDate(medallion string, pickUpDate time.Time, byPassCache bool) *mock.Call {
	return m.On("TripsByMedallionsOnPickUpDate", mock.AnythingOfType("*context.valueCtx"), medallion, pickUpDate, byPassCache)
}

func (m *mockTripSvc) TripsByMedallion(ctx context.Context, medallions []string, byPassCache bool) ([]output.Result, error) {
	args := m.Called(ctx, medallions, byPassCache)
	return args.Get(0).([]output.Result), args.Error(1)
}

func (m *mockTripSvc) OnTripsTripsByMedallion(medallions []string, byPassCache bool) *mock.Call {
	return m.On("TripsByMedallion", mock.AnythingOfType("*context.valueCtx"), medallions, byPassCache)
}
