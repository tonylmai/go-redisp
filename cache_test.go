package main

import (
	"testing"
	"time"
)

func TestManagedCache_Add(t *testing.T) {
	// Given
	cache = NewManagedCache(2, 5)

	// When we add more entries than capacity that one would get dropped from the list
	cache.Add("abc", "One")
	cache.Add("def", "Two")
	cache.Add("ghi", "Three")

	// Then
	if len(cache.keyAccessValues) != 2 {
		t.Errorf("Expected 2, got %d", len(cache.keyAccessValues))
	}
}

func TestManagedCache_Get(t *testing.T) {
	// Given
	cache = NewManagedCache(2, 5)
	cache.Add("abc", "Hello")

	// When/Then we access a cache, it should return the correct value
	var value = cache.Get("abc")
	if *value != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", value)
	}

	// When/Then we access a cache, it should increment the accessed count
	cache.Get("abc")
	cache.Get("abc")
	for _, accessValue := range cache.keyAccessValues {
		if accessValue.value == "Hello" {
			if accessValue.accessed != 3 {
				t.Errorf("Expected accessed count is 3, got %d", accessValue.accessed)
			}
		}
	}
}

func TestManagedCache_Capacity(t *testing.T) {
	// Given
	cache = NewManagedCache(2, 5)
	cache.Add("abc", "Hello")
	cache.Add("def", "World")

	// When I accessed it multiple times, the access counts would increment accordingly
	cache.Get("abc")
	cache.Get("abc")
	cache.Get("def")

	// Then when I add an entry that exceeds the capacity, the LRU would be dropped
	cache.Add("xyz", "Goodbye")
	var dropped = cache.Get("def")
	if dropped != nil {
		t.Errorf("Expected nil, got %s", dropped)
	}
}

func TestManagedCache_Expiry(t *testing.T) {
	// Given
	cache = NewManagedCache(2, 2)
	cache.Add("abc", "Hello")

	// When time has elapsed
	time.Sleep(3 * time.Second)

	// Then the cache should be dropped
	var dropped = cache.Get("abc")
	if dropped != nil {
		t.Errorf("Expected nil, got %s", dropped)
	}
}