package cache

import (
	"context"

	"github.com/muesli/cache2go"
)


// Cache embeds in memory cache2go.
type Cache struct {
	cache *cache2go.CacheTable 
}

// New creates a new instance of cache.
func New(cache *cache2go.CacheTable ) *Cache { return &Cache{cache:cache}}


// Get retrieves cache entries.
func (c *Cache) Get(ctx context.Context, key string)(int ,error) {
	res, err := c.cache.Value(key)
	if err != nil {
		return 0,err
	}
	return res.Data().(int),nil
}

// Set adds new cache entries.
func (c *Cache) Set(ctx context.Context,  key string, val int) {
	c.cache.Add(key, 0, val)

}

// Clear flush cache entries.
func (c *Cache) Clear(ctx context.Context) {
	c.cache.Flush()
}

