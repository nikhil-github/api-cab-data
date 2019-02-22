package service

import (
	"context"
	"time"
	"fmt"

	"github.com/nikhil-github/api-cab-data/pkg/output"
	"go.uber.org/zap"
)

const keyNotFound = "Key not found in cache"

type TripService struct {
	cacheGetter CacheGetter
	cacheSetter CacheSetter
	dbGetter Getter
	logger *zap.Logger
}

type Getter interface {
	Trips(ctx context.Context, medallion string,pickUpDate time.Time) (int, error)
}

type CacheGetter interface {
	Get(ctx context.Context,key string)(int ,error)
}

type CacheSetter interface {
	Set(ctx context.Context, key string, val int)
}

func New( g Getter,cg CacheGetter,cs CacheSetter,l *zap.Logger) *TripService {
	return &TripService{dbGetter: g,cacheGetter:cg,cacheSetter:cs,logger:l}
}

func (s *TripService) Trips(ctx context.Context, medallions []string,pickUpDate time.Time,byPassCache bool) ([]output.Result, error) {

	var results []output.Result
	for _, medallion := range medallions {
		result,err := s.get(ctx,medallion,pickUpDate,byPassCache)
		if err !=nil {
			s.logger.Error("Error finding trips for medallion:%s",zap.String("medallion",medallion))
			return []output.Result{},err
		}
		results = append(results, result)
	}
	return results,nil
}

func (s *TripService) get(ctx context.Context, medallion string,pickUpDate time.Time,byPassCache bool)(output.Result, error) {
	var count int
	var err error
	if byPassCache {
		count,err = s.getFromDB(ctx, medallion,pickUpDate)
		if err !=nil {
			return output.Result{},err
		}
	} else {
		count,err = s.cacheGetter.Get(ctx,key(medallion,pickUpDate))
		if err !=nil && err.Error() == keyNotFound {
			count,err = s.getFromDB(ctx, medallion,pickUpDate)
		}
	}
	return output.Result{Medallion:medallion,Trips:count},nil
}

func (s *TripService) getFromDB(ctx context.Context, medallion string,pickUpDate time.Time)(int ,error) {
	count,err := s.dbGetter.Trips(ctx,medallion,pickUpDate)
	if err !=nil {
		return 0,err
	}
	go s.cacheSetter.Set(ctx,key(medallion,pickUpDate),count)
	return count,nil
}


// key is built by concatinate cabID + pickUpDate.
func key(medallion string,pickUpDate time.Time,) string {
	return fmt.Sprintf("%s%d%d%d",medallion,pickUpDate.Year(),pickUpDate.Month(),pickUpDate.Day())
}