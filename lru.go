package gcache

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"
)

type lruValue struct {
	Key      interface{}
	Value    interface{}
	Deadline time.Time
}

func (v *lruValue) IsDeleted() bool {
	return !v.Deadline.IsZero() &&
		!v.Deadline.After(time.Now())
}

type LRU struct {
	opts lruOptions

	keys map[interface{}]*list.Element
	hot  *list.List

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
		opts:   opts,
		keys:   make(map[interface{}]*list.Element, opts.capacity),
		hot:    list.New(),
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
			l.unsafeClear()
			l.m.Unlock()
		}
	}
}
func (l *LRU) unsafeClear() {
	var (
		ele  *list.Element
		v    *lruValue
		hot  = l.hot
		keys = l.keys
	)
	if hot == nil {
		return
	}
	for {
		ele = hot.Front()
		if ele == nil {
			break
		}
		v = ele.Value.(*lruValue)
		if v.IsDeleted() {
			l.hot.Remove(ele)
			delete(keys, v.Key)
		} else {
			break
		}
	}
}

// Add the value to the cache, only when the key does not exist
func (l *LRU) Add(key, value interface{}) (added bool, e error) {
	l.m.Lock()
	defer l.m.Unlock()
	if l.done != 0 {
		e = ErrAlreadyClosed
		return
	}
	ele, exists := l.keys[key]
	if exists {
		v := ele.Value.(*lruValue)
		if v.IsDeleted() {
			added = true
			v.Value = value
			l.moveHot(ele)
		}
	} else {
		added = true
		l.add(key, value)
	}
	return
}
func (l *LRU) add(key, value interface{}) {
	// capacity limit reached, pop front
	for l.hot.Len() >= l.opts.capacity {
		ele := l.hot.Front()
		v := ele.Value.(*lruValue)
		delete(l.keys, v.Key)
		l.hot.Remove(ele)
	}
	// new value
	v := &lruValue{
		Key:   key,
		Value: value,
	}
	if l.opts.expiry > 0 {
		v.Deadline = time.Now().Add(l.opts.expiry)
	}
	l.keys[key] = l.hot.PushBack(v)
	return
}
func (l *LRU) moveHot(ele *list.Element) {
	v := ele.Value.(*lruValue)
	l.hot.Remove(ele)
	if l.opts.expiry > 0 {
		v.Deadline = time.Now().Add(l.opts.expiry)
	}
	l.keys[v.Key] = l.hot.PushBack(v)
}

// Put key value to cache
func (l *LRU) Put(key, value interface{}) (added bool, e error) {
	l.m.Lock()
	defer l.m.Unlock()
	if l.done != 0 {
		e = ErrAlreadyClosed
		return
	}
	ele, exists := l.keys[key]
	if exists {
		// put
		v := ele.Value.(*lruValue)
		if v.IsDeleted() {
			added = true
		}
		v.Value = value
		// move hot
		l.moveHot(ele)
	} else {
		added = true
		l.add(key, value)
	}
	return
}

// Get return cache value, if not exists then return ErrNotExists
func (l *LRU) Get(key interface{}) (value interface{}, e error) {
	l.m.Lock()
	defer l.m.Unlock()
	if l.done != 0 {
		e = ErrAlreadyClosed
		return
	}
	ele, exists := l.keys[key]
	if !exists {
		e = ErrNotExists
		return
	}
	v := ele.Value.(*lruValue)
	if v.IsDeleted() {
		delete(l.keys, key)
		l.hot.Remove(ele)
		e = ErrNotExists
		return
	}
	value = v.Value

	// move hot
	l.moveHot(ele)
	return
}

// Delete key from cache
func (l *LRU) Delete(key ...interface{}) (changed int, e error) {
	l.m.Lock()
	if l.done == 0 {
		var (
			ele    *list.Element
			exists bool
		)
		for _, k := range key {
			ele, exists = l.keys[k]
			if exists {
				changed++
				delete(l.keys, k)
				l.hot.Remove(ele)
			}
		}
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
		count = l.hot.Len()
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
		l.hot.Init()
		for k := range l.keys {
			delete(l.keys, k)
		}
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
			l.keys = nil
			l.hot = nil
			return nil
		}
	}
	return ErrAlreadyClosed
}
