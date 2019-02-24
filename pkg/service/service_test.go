package service_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/nikhil-github/api-cab-data/pkg/output"
	"github.com/nikhil-github/api-cab-data/pkg/service"
)

func TestTripService(t *testing.T) {
	t.Parallel()
	pDate := time.Date(2013, 12, 31, 0, 1, 0, 0, time.UTC)
	type args struct {
		Medallions  []string
		PickUpDate  time.Time
		ByPassCache bool
	}
	type fields struct {
		MockOperations func(d *dbMock, cg *cacheGetMock, cs *cacheSetMock)
		CacheSet       bool
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
			Name: "Get from DB",
			Args: args{Medallions: []string{"med1"}, PickUpDate: pDate, ByPassCache: true},
			Fields: fields{
				MockOperations: func(d *dbMock, cg *cacheGetMock, cs *cacheSetMock) {
					d.OnTrips("med1", pDate).Return(5, nil).Once()
					cs.OnSet("med120131231", 5)
					cs.wg = sync.WaitGroup{}
					cs.wg.Add(1)
				},
				CacheSet: true,
			},
			Want: want{Result: []output.Result{{Medallion: "med1", Trips: 5}}},
		},
		{
			Name: "Get from Cache",
			Args: args{Medallions: []string{"med2"}, PickUpDate: pDate, ByPassCache: false},
			Fields: fields{
				MockOperations: func(d *dbMock, cg *cacheGetMock, cs *cacheSetMock) {
					cg.OnGet("med220131231").Return(10, nil).Once()
				},
			},
			Want: want{Result: []output.Result{{Medallion: "med2", Trips: 10}}},
		},
		{
			Name: "Cache Missed",
			Args: args{Medallions: []string{"med2"}, PickUpDate: pDate, ByPassCache: false},
			Fields: fields{
				MockOperations: func(d *dbMock, cg *cacheGetMock, cs *cacheSetMock) {
					cg.OnGet("med220131231").Return(0, errors.New("Key not found in cache"))
					d.OnTrips("med2", pDate).Return(5, nil).Once()
					cs.OnSet("med220131231", 5)
					cs.wg = sync.WaitGroup{}
					cs.wg.Add(1)
				},
				CacheSet: true,
			},
			Want: want{Result: []output.Result{{Medallion: "med2", Trips: 5}}},
		},
		{
			Name: "Failure",
			Args: args{Medallions: []string{"med3"}, PickUpDate: pDate, ByPassCache: true},
			Fields: fields{
				MockOperations: func(d *dbMock, cg *cacheGetMock, cs *cacheSetMock) {
					d.OnTrips("med3", pDate).Return(0, errors.New("error"))
				},
			},
			Want: want{Error: "error"},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.Name, func(t *testing.T) {
			var db dbMock
			var cacheGet cacheGetMock
			var cacheSet cacheSetMock
			tt.Fields.MockOperations(&db, &cacheGet, &cacheSet)
			svc := service.New(&db, &cacheGet, &cacheSet, zap.NewNop())
			result, err := svc.Trips(context.Background(), tt.Args.Medallions, tt.Args.PickUpDate, tt.Args.ByPassCache)
			if tt.Fields.CacheSet {
				cacheSet.wg.Wait()
			}
			if tt.Want.Error != "" {
				assert.EqualError(t, err, tt.Want.Error)
				return
			}
			require.NoError(t, err, "should not return an error")
			assert.Equal(t, tt.Want.Result, result, "results")
		})
	}

}

type dbMock struct {
	mock.Mock
}

func (d *dbMock) Trips(ctx context.Context, medallion string, pickUpDate time.Time) (int, error) {
	args := d.Called(ctx, medallion, pickUpDate)
	return args.Get(0).(int), args.Error(1)
}

func (d *dbMock) OnTrips(medallion string, pickUpDate time.Time) *mock.Call {
	return d.On("Trips", mock.AnythingOfTypeArgument("*context.emptyCtx"), medallion, pickUpDate)
}

type cacheGetMock struct {
	mock.Mock
}

func (cg *cacheGetMock) Get(ctx context.Context, key string) (int, error) {
	args := cg.Called(ctx, key)
	return args.Get(0).(int), args.Error(1)
}

func (cg *cacheGetMock) OnGet(key string) *mock.Call {
	return cg.On("Get", mock.AnythingOfTypeArgument("*context.emptyCtx"), key)
}

type cacheSetMock struct {
	mock.Mock
	wg sync.WaitGroup
}

func (cs *cacheSetMock) Set(ctx context.Context, key string, val int) {
	cs.Called(ctx, key, val)
	cs.wg.Done()
	return
}

func (cs *cacheSetMock) OnSet(key string, val int) *mock.Call {
	return cs.On("Set", mock.AnythingOfTypeArgument("*context.emptyCtx"), key, val)
}
