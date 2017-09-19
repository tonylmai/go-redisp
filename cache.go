package main

import (
	"sync"
	"time"
	"math"
	"log"
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

// Get value from cache. If cache has expired, delete it from the cache and return nil ("" is a valid value)
func (c *managedCache) Get(key string) *string {
	var accessVal =  c.keyAccessValues[key]
	// How long have you been in the cache?
	now := time.Now()
	diff := now.Second() - accessVal.time.Second()

	log.Printf("Now=%d, then=%d, diff=%d\n", now.Second(), accessVal.time.Second(), diff)

	if diff > int(c.expiry) {
		log.Printf("Cache is old. Removing key=%s", key)
		delete(c.keyAccessValues, key)
		return nil
	} else {
		// Increment by 1
		accessVal.accessed += 1
		log.Printf("For key=%s, setting access count to %d\n", key, accessVal.accessed)
		return &accessVal.value
	}
}

// Add to cache. Should cache is at capacity, delete the least used entry before adding the new entry
func (c *managedCache) Add(key string, value string) {
	log.Printf("Adding %s=%s\n", key, value)
	// If at capacity, drop the LRU in the map (sorted by accessed)
	if int64(len(c.keyAccessValues)) == c.capacity {
		var lru = getLRUKey(c.keyAccessValues)
		log.Printf("cach is at Capacity. Removing key=%s\n", lru)
		delete(c.keyAccessValues, lru)
	}

	// Now add to map
	log.Printf("Adding new key %s=%s\n", key, value)
	t := time.Now()
	c.keyAccessValues[key] = accessValue{t, 1, value,}
}

// Get the LRU entry in the cache. Should there be more than one, chose the first encounter
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