package cabs

import (
	"net/http"
	"github.com/go-chi/chi"
	"strings"
	"encoding/json"

	"log"
	"time"
	"fmt"
	"context"
)

type Handler struct{
	svc Servicer
}


type Servicer interface {
	Trips(ctx context.Context, cabIDs []string,pickUpDate time.Time,byPassCache bool) ([]Result, error)
}


func NewHandler(tripSvc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		pickUpDate := chi.URLParam(r, "pickupDate")

		d, err := time.Parse("2006-01-02", pickUpDate)
		if err != nil {
			log.Printf("error: pickUpDate is not a valid date")
			responseBadRequest(w,enc,"invalid pick up date")
			return
		}
		log.Println(d)

		medallions := strings.Split(chi.URLParam(r, "ids"), ",")
		fmt.Println(err)
		if len(medallions) == 0 {
			log.Println("invalid medallions")
			responseBadRequest(w,enc,"invalid medallions")
		}
		log.Println("medallions",medallions)
		log.Println("medallions length",len(medallions))
		fmt.Println("pik date",d.String())
		tripSvc.Trips(r.Context(),medallions,d,false)

		w.Write([]byte("OK"))
	}
}
//
//
//
//func Handle(w http.ResponseWriter, r *http.Request) {
//	// ctx := r.Context()
//	enc := json.NewEncoder(w)
//	pickUpDate := chi.URLParam(r, "pickupDate")
//
//	d, err := time.Parse("2006-01-02", pickUpDate)
//	if err != nil {
//		log.Printf("error: pickUpDate is not a valid date")
//		responseBadRequest(w,enc,"invalid pick up date")
//		return
//	}
//	log.Println(d)
//
//	medallions := strings.Split(chi.URLParam(r, "ids"), ",")
//	fmt.Println(err)
//	if len(medallions) == 0 {
//		log.Println("invalid medallions")
//		responseBadRequest(w,enc,"invalid medallions")
//	}
//	log.Println("medallions",medallions)
//	log.Println("medallions length",len(medallions))
//
//	Servicer
//
//	w.Write([]byte("OK"))
//}

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
