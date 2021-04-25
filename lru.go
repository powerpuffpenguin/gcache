package gcache

import (
	"sync"
	"sync/atomic"
	"time"
)

type LRU struct {
	opts *lruOptions

	impl *lru

	ticker *time.Ticker

	closed chan struct{}
	done   uint32
	m      sync.Mutex
}

func NewLRU(opt ...LRUOption) (lru *LRU) {
	opts := defaultLRUOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	lru = &LRU{
		opts:   &opts,
		impl:   newLRU(&opts),
		closed: make(chan struct{}),
	}
	if opts.expiry > 0 {
		ticker := time.NewTicker(opts.clear)
		lru.ticker = ticker
		go lru.clear(ticker.C)
	}
	return
}

func (l *LRU) clear(ch <-chan time.Time) {
	for {
		select {
		case <-l.closed:
			return
		case <-ch:
			l.m.Lock()
			if l.done == 0 {
				l.impl.ClearExpired()
			}
			l.m.Unlock()
		}
	}
}

// Add the value to the cache, only when the key does not exist
func (l *LRU) Add(key, value interface{}) (newkey bool, e error) {
	l.m.Lock()
	if l.done == 0 {
		newkey = l.impl.Add(key, value)
	} else {
		e = ErrAlreadyClosed
	}
	l.m.Unlock()
	return
}

// Put key value to cache
func (l *LRU) Put(key, value interface{}) (newkey bool, e error) {
	l.m.Lock()
	if l.done == 0 {
		newkey = l.impl.Put(key, value)
	} else {
		e = ErrAlreadyClosed
	}
	l.m.Unlock()
	return
}

// Get return cache value, if not exists then return ErrNotExists
func (l *LRU) Get(key interface{}) (value interface{}, e error) {
	l.m.Lock()
	if l.done == 0 {
		value, e = l.impl.Get(key)
	} else {
		e = ErrAlreadyClosed
	}
	l.m.Unlock()
	return
}

// Delete key from cache
func (l *LRU) Delete(key ...interface{}) (changed int, e error) {
	l.m.Lock()
	if l.done == 0 {
		changed = l.impl.Delete(key...)
	} else {
		e = ErrAlreadyClosed
	}
	l.m.Unlock()
	return
}

// Len returns the number of cached data
func (l *LRU) Len() (count int, e error) {
	l.m.Lock()
	if l.done == 0 {
		count = l.impl.Len()
	} else {
		e = ErrAlreadyClosed
	}
	l.m.Unlock()
	return
}

// Clear all cached data
func (l *LRU) Clear() (e error) {
	l.m.Lock()
	if l.done == 0 {
		l.Clear()
	} else {
		e = ErrAlreadyClosed
	}
	l.m.Unlock()
	return
}

// Close cache
func (l *LRU) Close() (e error) {
	if atomic.LoadUint32(&l.done) == 0 {
		l.m.Lock()
		defer l.m.Unlock()
		if l.done == 0 {
			defer atomic.StoreUint32(&l.done, 1)
			close(l.closed)
			if l.ticker != nil {
				l.ticker.Stop()
			}
			l.impl.Clear()
			return nil
		}
	}
	return ErrAlreadyClosed
}
