package gcache

import "errors"

var (
	ErrNotExists     = errors.New(`key not exists`)
	ErrAlreadyClosed = errors.New(`cache already closed`)
)

type Value struct {
	Exists bool
	Value  interface{}
}

type Cache interface {
	// Add the value to the cache, only when the key does not exist
	Add(key, value interface{}) (added bool)
	// Put key value to cache
	Put(key, value interface{}) (added bool)
	// Get return cache value, if not exists then return ErrNotExists
	Get(key interface{}) (value interface{})
	// BatchPut pairs to cache
	BatchPut(pair ...interface{})
	// BatchGet return cache values
	BatchGet(key ...interface{}) (vals []Value)
	// Delete key from cache
	Delete(key ...interface{}) (changed int)
	// Len returns the number of cached data
	Len() (count int)
	// Clear all cached data
	Clear()
}

// Low-level caching is usually only used when combining multiple caching algorithms
type LowCache interface {
	// Clear Expired cache
	ClearExpired()
	// Add the value to the cache, only when the key does not exist
	Add(key, value interface{}) (added bool)
	// Put key value to cache
	Put(key, value interface{}) (added bool)
	// Get return cache value
	Get(key interface{}) (value interface{}, exists bool)
	// Delete key from cache
	Delete(key ...interface{}) (changed int)
	// Len returns the number of cached data
	Len() int
	// Clear all cached data
	Clear()
}
