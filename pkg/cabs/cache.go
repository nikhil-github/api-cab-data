package cabs

import (
	"context"
	"fmt"
	"time"

	"github.com/muesli/cache2go"
)


type Cache struct {
	cache *cache2go.CacheTable 
}

func NewCache(cache *cache2go.CacheTable ) *Cache { return &Cache{cache:cache}}


func (c *Cache) Get(ctx context.Context, cabID string,pickUpDate time.Time)(int ,error) {
	key := key(cabID,pickUpDate)
	res, err := c.cache.Value(key)
	fmt.Println("cache key is ",key, res)
	fmt.Println("cache value is ",err)

	if err != nil {
		fmt.Println("get cahce",err)
		return 0,err
	}
	fmt.Println("value from cache ",key,res.Data().(int))
	return res.Data().(int),nil
}

func (c *Cache) Set(ctx context.Context, cabID string,pickUpDate time.Time, count int) {
	key := key(cabID,pickUpDate)
	fmt.Println("set cache key is ",key)
	c.cache.Add(key, 0, count)

}

func (c *Cache) Clear(ctx context.Context) {
	c.cache.Flush()
}

func key(cabID string,pickUpDate time.Time,) string {
	return fmt.Sprintf("%s%d%d%d",cabID,pickUpDate.Year(),pickUpDate.Month(),pickUpDate.Day())
}