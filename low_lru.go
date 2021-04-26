package gcache

import (
	"container/list"
	"time"
)

// A low-level implementation of lru, use LRU unless you know exactly what you are doing.
type LowLRU struct {
	opts lowLRUOptions
	keys map[interface{}]*list.Element
	hot  *list.List
}

// NewLowLRU create a low-level lru, use NewLRU unless you know exactly what you are doing.
func NewLowLRU(opt ...LowLRUOption) *LowLRU {
	opts := defaultLowLRUOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	return &LowLRU{
		opts: opts,
		keys: make(map[interface{}]*list.Element, opts.capacity),
		hot:  list.New(),
	}
}
func (l *LowLRU) ClearExpired() {
	var (
		ele  *list.Element
		v    cacheValue
		hot  = l.hot
		keys = l.keys
	)
	for {
		ele = hot.Front()
		if ele == nil {
			break
		}
		v = ele.Value.(cacheValue)
		if v.IsDeleted() {
			l.hot.Remove(ele)
			delete(keys, v.GetKey())
		} else {
			break
		}
	}
}

// Add the value to the cache, only when the key does not exist
func (l *LowLRU) Add(key, value interface{}) (added bool) {
	ele, exists := l.keys[key]
	if exists {
		v := ele.Value.(cacheValue)
		if v.IsDeleted() {
			added = true
			v.SetValue(value)
			l.moveHot(ele)
		}
	} else {
		added = true
		l.add(key, value)
	}
	return
}
func (l *LowLRU) add(key, value interface{}) (delkey, delval interface{}, deleted bool) {
	// capacity limit reached, pop front
	if l.hot.Len() >= l.opts.capacity {
		deleted = true
		ele := l.hot.Front()
		v := ele.Value.(cacheValue)
		delkey = v.GetKey()
		delval = v.GetValue()
		delete(l.keys, delkey)
		l.hot.Remove(ele)
	}
	// new value
	v := newValue(key, value, l.opts.expiry)
	l.keys[key] = l.hot.PushBack(v)
	return
}
func (l *LowLRU) moveHot(ele *list.Element) {
	v := ele.Value.(cacheValue)
	l.hot.Remove(ele)
	if l.opts.expiry > 0 {
		v.SetDeadline(time.Now().Add(l.opts.expiry))
	}
	l.keys[v.GetKey()] = l.hot.PushBack(v)
}

func (l *LowLRU) Put(key, value interface{}) (delkey, delval interface{}, deleted bool) {
	ele, exists := l.keys[key]
	if exists {
		// put
		v := ele.Value.(cacheValue)
		v.SetKey(value)
		// move hot
		l.moveHot(ele)
	} else {
		delkey, delval, deleted = l.add(key, value)
	}
	return
}

// Get return cache value
func (l *LowLRU) Get(key interface{}) (value interface{}, exists bool) {
	ele, exists := l.keys[key]
	if !exists {
		return
	}
	v := ele.Value.(cacheValue)
	if v.IsDeleted() {
		delete(l.keys, key)
		l.hot.Remove(ele)
		exists = false
		return
	}
	value = v.GetValue()

	// move hot
	l.moveHot(ele)
	return
}
func (l *LowLRU) Delete(key ...interface{}) (changed int) {
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
	return
}

func (l *LowLRU) Len() int {
	return l.hot.Len()
}

func (l *LowLRU) Clear() {
	l.hot.Init()
	for k := range l.keys {
		delete(l.keys, k)
	}
	return
}
