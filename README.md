# gcache

golang cache interface and some algorithm implementation

* fifo
* lfu
* lru
* lru-k
* 2q

# example 

```
github.com/powerpuffpenguin/gcache
```

```
package main

import (
	"fmt"

	"github.com/powerpuffpenguin/gcache"
)

func main() {
	c := gcache.NewLRU(
		gcache.WithLRUCapacity(3),
	)
	for i := 0; i < 4; i++ {
		key := i
		val := fmt.Sprintf(`val %d`, i)
		c.Put(key, val)
	}

	for i := 0; i < 4; i++ {
		key := i
		val, exists := c.Get(key)
		if exists {
			fmt.Printf("key=%v val=%v\n", key, val)
		} else {
			fmt.Printf("key=%v not exists\n", key)
		}
	}
}
```

# interface 

gcache provides two interface for users to use.

1. **Cache** goroutine safe cache
2. **LowCache** Implementation of goroutine not safe cache low-level algorithm

## Cache

Usually you should use the Cache interface directly.

```
type Cache interface {
	// Add the value to the cache, only when the key does not exist
	Add(key, value interface{}) (added bool)
	// Put key value to cache
	Put(key, value interface{})
	// Get return cache value
	Get(key interface{}) (value interface{}, exists bool)
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
```

gcache provides several NewXXX functions for creating Cache. In addition, several WithXXX settings are provided to guide parameters for the caching algorithm.

```
capacity := 1000
expiry := time.Minute
duration := time.Minute
// basic
lru := gcache.NewLRU(
	gcache.WithLRUCapacity(capacity), // cache capacity
	gcache.WithLRUExpiry(expiry),     // inactivity expiration time
	gcache.WithLRUClear(duration),    // the timer clears the expired cache
)
lfu := gcache.NewLFU(
	gcache.WithLFUCapacity(capacity),
	gcache.WithLFUExpiry(expiry),
	gcache.WithLFUClear(duration),
)
fifo := gcache.NewFIFO(
// WithFIFOXXX
)
// complex
lruk2 := gcache.NewLRUK(
	gcache.WithLRUK(2),
	gcache.WithLRUKHistoryOnlyKey(true),        // if true history only save key not save value, false save key and value
	gcache.WithLRUKHistory(gcache.NewLowLRU()), // history use lru
)
lruk3 := gcache.NewLRUK(
	gcache.WithLRUK(3),
)
// 2q
c2q := gcache.NewLRUK(
	gcache.WithLRUK(2),
	gcache.WithLRUKHistoryOnlyKey(false),
	gcache.WithLRUKHistory(gcache.NewLowFIFO()), // history use fifo
)
```

## LowCache

The LowCache interface is a low-level implementation that implements the basic algorithm.

```
// Low-level caching is usually only used when combining multiple caching algorithms
type LowCache interface {
	// Clear Expired cache
	ClearExpired()
	// Add the value to the cache, only when the key does not exist
	Add(key, value interface{}) (added bool)
	// Put key value to cache
	Put(key, value interface{}) (delkey, delval interface{}, deleted bool)
	// Get return cache value
	Get(key interface{}) (value interface{}, exists bool)
	// Delete key from cache
	Delete(key ...interface{}) (changed int)
	// Len returns the number of cached data
	Len() int
	// Clear all cached data
	Clear()
}
```

gcache provides several NewLowXXX functions for creating LowCache. In addition, several WithLowXXX settings are provided to guide parameters for the caching algorithm.

# fifo

```
gcache.NewFIFO(
	gcache.WithFIFOCapacity(1000),
	gcache.WithFIFOExpiry(time.Minute),
	gcache.WithFIFOClear(time.Minute*10),
)
```
# lfu
```
gcache.NewLFU(
	gcache.WithLFUCapacity(1000),
	gcache.WithLFUExpiry(time.Minute),
	gcache.WithLFUClear(time.Minute*10),
)
```

# lru

```
gcache.NewLRU(
	gcache.WithLRUCapacity(1000),
	gcache.WithLRUExpiry(time.Minute), 
	gcache.WithLRUClear(time.Minute*10),
)
```

# lru-k

```
gcache.NewLRUK(
	gcache.WithLRUKCapacity(1000),
	gcache.WithLRUKExpiry(time.Minute),
	gcache.WithLRUKClear(time.Minute*10),

	gcache.WithLRUK(2),
	gcache.WithLRUKHistoryOnlyKey(false),
	gcache.WithLRUKHistory(gcache.NewLowLRU()),
)
```

# 2q

```
gcache.NewLRUK(
	gcache.WithLRUKCapacity(1000),
	gcache.WithLRUKExpiry(time.Minute),
	gcache.WithLRUKClear(time.Minute*10),

	gcache.WithLRUK(2),
	gcache.WithLRUKHistoryOnlyKey(false),
	gcache.WithLRUKHistory(gcache.NewLowFIFO(
		gcache.WithLowFIFOCapacity(1000),
		gcache.WithLowFIFOExpiry(time.Minute),
	)),
)
```