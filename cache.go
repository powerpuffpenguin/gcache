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
	Add(key, value interface{}) (newkey bool, e error)
	// Put key value to cache
	Put(key, value interface{}) (newkey bool, e error)
	// Get return cache value, if not exists then return ErrNotExists
	Get(key interface{}) (value interface{}, e error)
	// BatchPut pairs to cache
	BatchPut(pair ...interface{}) (e error)
	// BatchGet return cache values
	BatchGet(key ...interface{}) (vals []Value, e error)
	// Delete key from cache
	Delete(key ...interface{}) (changed int, e error)
	// Len returns the number of cached data
	Len() (count int, e error)
	// Clear all cached data
	Clear() (e error)
	// Close cache
	Close() (e error)
}
