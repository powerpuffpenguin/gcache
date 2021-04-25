package gcache

import (
	"sync"
	"time"
)

type wrapper struct {
	impl   LowCache
	ticker *time.Ticker

	closed chan struct{}
	m      sync.Mutex
}

// Add the value to the cache, only when the key does not exist
func (w *wrapper) Add(key, value interface{}) (added bool) {
	w.m.Lock()
	added = w.impl.Add(key, value)
	w.m.Unlock()
	return
}

// Put key value to cache
func (w *wrapper) Put(key, value interface{}) {
	w.m.Lock()
	w.impl.Put(key, value)
	w.m.Unlock()
	return
}

// Get return cache value, if not exists then return ErrNotExists
func (w *wrapper) Get(key interface{}) (value interface{}, exists bool) {
	w.m.Lock()
	value, exists = w.impl.Get(key)
	w.m.Unlock()
	return
}

// BatchPut pairs to cache
func (w *wrapper) BatchPut(pair ...interface{}) {
	w.m.Lock()
	count := len(pair)
	for i := 0; i < count; i += 2 {
		if i+1 < count {
			w.impl.Put(pair[i], pair[i+1])
		} else {
			w.impl.Put(pair[i], nil)
			break
		}
	}
	w.m.Unlock()
	return
}

// BatchGet return cache values
func (w *wrapper) BatchGet(key ...interface{}) (vals []Value) {
	w.m.Lock()
	vals = make([]Value, len(key))
	for i, k := range key {
		vals[i].Value, vals[i].Exists = w.impl.Get(k)
	}
	w.m.Unlock()
	return
}

// Delete key from cache
func (w *wrapper) Delete(key ...interface{}) (changed int) {
	w.m.Lock()
	changed = w.impl.Delete(key...)
	w.m.Unlock()
	return
}

// Len returns the number of cached data
func (w *wrapper) Len() (count int) {
	w.m.Lock()
	count = w.impl.Len()
	w.m.Unlock()
	return
}

// Clear all cached data
func (w *wrapper) Clear() {
	w.m.Lock()
	w.impl.Clear()
	w.m.Unlock()
}

func (w *wrapper) clearExpired(ch <-chan time.Time) {
	for {
		select {
		case <-w.closed:
			return
		case <-ch:
			w.m.Lock()
			w.impl.ClearExpired()
			w.m.Unlock()
		}
	}
}
func (w *wrapper) stop() {
	w.m.Lock()
	close(w.closed)
	if w.ticker != nil {
		w.ticker.Stop()
	}
	w.impl.Clear()
	w.m.Unlock()
}
