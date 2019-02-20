package cabs

import (
	"context"
	"time"
	"fmt"
)


type Result struct{
	CabID string
	Count int
}

type Service struct {
	cacheGetter CacheGetter
	cacheSetter CacheSetter
	dbGetter Getter
}

type Getter interface {
	CabTripsByPickUpDate(ctx context.Context, medallion string,pickUpDate time.Time) (int, error)
}

type CacheGetter interface {
	Get(ctx context.Context, cabID string,pickUpDate time.Time)(int ,error)
}

type CacheSetter interface {
	Set(ctx context.Context, cabID string,pickUpDate time.Time, count int)
}

func NewService( g Getter,cg CacheGetter,cs CacheSetter) *Service {
	return &Service{dbGetter:g,cacheGetter:cg,cacheSetter:cs}
}

func (c *Service) Trips(ctx context.Context, cabIDs []string,pickUpDate time.Time,byPassCache bool) ([]Result, error) {

	var result []Result
	for _,cabID := range cabIDs {
		var count int
		var err error
		if byPassCache {
			count,err = c.getFromDB(ctx,cabID,pickUpDate)
		} else {
			count,err = c.cacheGetter.Get(ctx,cabID,pickUpDate)
			fmt.Println("count from cache",count)
			fmt.Println("error from cache",err)
			if err !=nil && err.Error() == "Key not found in cache" {
				count,err = c.getFromDB(ctx,cabID,pickUpDate)
			}
		}

		result = append(result,Result{cabID,count})
	}

	fmt.Println(result)
	return result,nil
}

func (c *Service) getFromDB(ctx context.Context, cabID string,pickUpDate time.Time)(int ,error) {
	count,_ := c.dbGetter.CabTripsByPickUpDate(ctx,cabID,pickUpDate)
	go c.cacheSetter.Set(ctx,cabID,pickUpDate,count)
	return count,nil
}