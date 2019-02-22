package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/nikhil-github/api-cab-data/pkg/output"
)

type Params struct {
	Logger *zap.Logger
	Svc    TripServicer
	Cache  Clearer
}

type TripServicer interface {
	Trips(ctx context.Context, cabIDs []string,pickUpDate time.Time,byPassCache bool) ([]output.Result, error)
}

type Clearer interface {
	Clear(ctx context.Context)
}

// Trips query for number of trips per medallion.
func Trips(logger *zap.Logger,tripSvc TripServicer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		pickUpDateStr := mux.Vars(r)["pickupDate"]
		pickupDate, err := time.Parse("2006-01-02", pickUpDateStr)
		if err != nil {
			logger.Error("error: pickUpDate is not a valid date",zap.String("pickUpDateStr",pickUpDateStr))
			responseBadRequest(w,enc,"invalid pick up date")
			return
		}

		cabIDs := strings.Split( mux.Vars(r)["ids"], ",")
		if len(cabIDs) == 0 {
			logger.Error("medallions missing")
			responseBadRequest(w,enc,"invalid medallions")
		}

		byPassCache,err := parseQueryParam(r)
		if err != nil {
			logger.Error("ByPassCache is not a valid",zap.Bool("byPassCache",byPassCache))
			responseBadRequest(w,enc,"invalid bypasscache")
			return
		}

		results,err := tripSvc.Trips(r.Context(),cabIDs,pickupDate,byPassCache)
		if err != nil {
			logger.Error("Error: counting trips",zap.Error(err))
			serverError(w,enc,"service failure")
			return
		}
		responseOK(w,enc,results)
	}
}

// ClearCache flushes the cache entries.
func ClearCache(logger *zap.Logger,cache Clearer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cache.Clear(r.Context())
		logger.Info("flushed cache entries")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`OK`))
	}
}

func parseQueryParam(r *http.Request) (bool,error) {
	queryValues := r.URL.Query()
	val := queryValues.Get("bypasscache")
	if len(val) == 0 {
		return false,nil
	}
	flag, err := strconv.ParseBool(val)
	if err!= nil {
		return false,err
	}
	return flag,nil
}

func responseOK(w http.ResponseWriter, encoder *json.Encoder, response interface{}) {
	w.WriteHeader(http.StatusOK)
	encoder.Encode(response)
}

func responseBadRequest(w http.ResponseWriter, encoder *json.Encoder, response string) {
	w.WriteHeader(http.StatusBadRequest)
	encoder.Encode(NewErrorMsg(response))
}

func serverError(w http.ResponseWriter, encoder *json.Encoder,response string) {
	code := http.StatusInternalServerError
	w.WriteHeader(code)
	encoder.Encode(NewErrorMsg(response))
}

//ErrorMsg construct for application error
type ErrorMsg struct {
	Message    string `json:"message"`
}

//NewErrorMsg returns ErrorMsg
func NewErrorMsg(message string) *ErrorMsg {
	return &ErrorMsg{
		Message:    message,
	}
}
