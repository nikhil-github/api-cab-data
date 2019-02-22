package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"time"
	"github.com/nikhil-github/api-cab-data/pkg/service"
	"github.com/go-chi/chi"
)

func TestHandler(t *testing.T) {
	type args struct {
		URL string
		Path string
	}
	type fields struct {
		MockExpectations func(m *mockTripSvc)
	}
	type want struct {
		Status int
	}
	testTable := []struct {
		Name   string
		Args   args
		Fields fields
		Want   want
	}{
		{
			Name:   "Failure - Invalid pickup date",
			Args:   args{Path:"/trips/pickupdate/2013-12-31/medallion/67EB082BFFE72095EAF18488BEA96050"},
			Fields: fields{MockExpectations: func(m *mockTripSvc) {}},
			Want:   want{Status: http.StatusBadRequest},
		},
	}
	for _, tt := range testTable {
		t.Run(tt.Name, func(t *testing.T) {
			logger := zap.NewNop()
			var m mockTripSvc
			tt.Fields.MockExpectations(&m)
			w := httptest.NewRecorder()
			r ,err := http.NewRequest("GET", tt.Args.Path, nil)
			rctx := chi.NewRouteContext()
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			require.NoError(t, err, "failed to create http GET request")
			h := service.NewHandler(logger, &m)
			h(w, r)
			m.AssertExpectations(t)
			assert.Equal(t, tt.Want.Status, w.Code, "status")
		})
	}
}

type mockTripSvc struct {
	mock.Mock
}

func (m *mockTripSvc) Trips(ctx context.Context, cabIDs []string,pickUpDate time.Time,byPassCache bool) ([]service.Result, error){
	args := m.Called(ctx, cabIDs,pickUpDate,byPassCache)
	return args.Get(0).([]service.Result),args.Error(0)
}
