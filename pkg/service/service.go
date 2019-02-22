package service

import (
	"context"
	"time"
	"fmt"

	"github.com/nikhil-github/api-cab-data/pkg/output"
)

type Service struct {
	cacheGetter CacheGetter
	cacheSetter CacheSetter
	dbGetter Getter
}

type Getter interface {
	TripsByPickUpDate(ctx context.Context, medallion string,pickUpDate time.Time) (int, error)
}

type CacheGetter interface {
	Get(ctx context.Context,key string)(int ,error)
}

type CacheSetter interface {
	Set(ctx context.Context, key string, val int)
}

func NewService( g Getter,cg CacheGetter,cs CacheSetter) *Service {
	return &Service{dbGetter:g,cacheGetter:cg,cacheSetter:cs}
}

func (s *Service) Trips(ctx context.Context, medallions []string,pickUpDate time.Time,byPassCache bool) ([]output.Result, error) {

	var result []output.Result
	for _, medallion := range medallions {
		var count int
		var err error
		if byPassCache {
			count,err = s.getFromDB(ctx, medallion,pickUpDate)
		} else {
			count,err = s.cacheGetter.Get(ctx,key(medallion,pickUpDate))
			if err !=nil && err.Error() == "Key not found in cache" {
				count,err = s.getFromDB(ctx, medallion,pickUpDate)
			}
		}

		result = append(result, output.Result{medallion,count})
	}
	return result,nil
}

func (s *Service) getFromDB(ctx context.Context, cabID string,pickUpDate time.Time)(int ,error) {
	count,err := s.dbGetter.TripsByPickUpDate(ctx,cabID,pickUpDate)
	if err !=nil {
		return 0,err
	}
	go s.cacheSetter.Set(ctx,key(cabID,pickUpDate),count)
	return count,nil
}


// key is built by concatinate cabID + pickUpDate.
func key(cabID string,pickUpDate time.Time,) string {
	return fmt.Sprintf("%s%d%d%d",cabID,pickUpDate.Year(),pickUpDate.Month(),pickUpDate.Day())
}