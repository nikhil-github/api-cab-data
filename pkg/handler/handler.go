package handler

import (
	"net/http"
	"github.com/go-chi/chi"
	"strings"
	"encoding/json"

	"time"
	"context"
	"go.uber.org/zap"
	"strconv"
	"github.com/nikhil-github/api-cab-data/pkg/service"
)


type TripServicer interface {
	Trips(ctx context.Context, cabIDs []string,pickUpDate time.Time,byPassCache bool) ([]service.Result, error)
}

type Clearer interface {
	Clear(ctx context.Context)
}


func Trips(logger *zap.Logger,tripSvc TripServicer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		pickUpDateStr := chi.URLParam(r, "pickupDate")

		pickupDate, err := time.Parse("2006-01-02", pickUpDateStr)
		if err != nil {
			logger.Error("error: pickUpDate is not a valid date",zap.String("pickUpDateStr",pickUpDateStr))
			responseBadRequest(w,enc,"invalid pick up date")
			return
		}

		cabIDs := strings.Split(chi.URLParam(r, "ids"), ",")
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
			serverError(w,enc)
			return
		}
		responseOK(w,enc,results)
	}
}

func ClearCache(logger *zap.Logger,cache Clearer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cache.Clear(r.Context())
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
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
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	encoder.Encode(response)
}

func responseBadRequest(w http.ResponseWriter, encoder *json.Encoder, response interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
	encoder.Encode(response)
}

func serverError(w http.ResponseWriter, encoder *json.Encoder) {
	code := http.StatusInternalServerError
	w.WriteHeader(code)
	// encoder.Encode(data.GeneralError(http.StatusText(code)))
}

//AppError construct for application error
type AppError struct {
	Error      error  `json:"error"`
	StatusCode int    `json:"httpStatusCode"`
	Message    string `json:"Message"`
}

//NewAppError returns AppError
func NewAppError(err error, statusCode int, message string) *AppError {
	return &AppError{
		Error:      err,
		StatusCode: statusCode,
		Message:    message,
	}
}
