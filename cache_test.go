package main

import "testing"

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