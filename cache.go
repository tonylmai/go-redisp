package main

import (
	"sync"
	"time"
	"math"
)

// Managed cache with a lock
type managedCache struct {
	sync.Mutex
	capacity        int64
	expiry          int64
	keyAccessValues map[string]accessValue
}

// A accessed/value pair
type accessValue struct {
	time     time.Time
	accessed int
	value    string
}

// Create a new ManagedCache struct
func NewManagedCache(capacity int64, expiry int64) *managedCache {
	m := managedCache{
		capacity:        capacity,
		expiry:          expiry,
		keyAccessValues: make(map[string]accessValue, capacity),
	}
	return &m
}

// Get value from cache. If cache has expired, delete it from the cache and return nil
func (c *managedCache) Get(key string) *string {
	var accessVal =  c.keyAccessValues[key]
	// How long have you been in the cache?
	diff := time.Now().Sub(accessVal.time)
	if int64(diff.Seconds()) > c.expiry {
		delete(c.keyAccessValues, key)
		return nil
	} else {
		// Increment by 1
		accessVal.accessed += 1
		return &accessVal.value
	}
}

// Add to cache. Should cache is at capacity, delete the least used entry before adding the new entry
func (c *managedCache) Add(key string, value string) {
	// If at capacity, drop the LRU in the map (sorted by accessed)
	if int64(len(c.keyAccessValues)) == c.capacity {
		var key = getLRUKey(c.keyAccessValues)
		delete(c.keyAccessValues, key)
	}

	// Now add to both maps
	t := time.Now()
	c.keyAccessValues[key] = accessValue{t, 1, value,}
}

// Get the oldest entry in the cache
func getLRUKey(m map[string]accessValue) string {
	frequency := math.MaxInt64
	lru := ""
	for k, v  := range m {
		if v.accessed < frequency {
			frequency = v.accessed
			lru = k
		}
	}
	return lru
}