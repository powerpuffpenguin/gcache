package gcache

import (
	"container/list"
	"time"
)

type cacheValue struct {
	Key      interface{}
	Value    interface{}
	Deadline time.Time
}

func (v *cacheValue) IsDeleted() bool {
	return !v.Deadline.IsZero() &&
		!v.Deadline.After(time.Now())
}

// reserved internal lru algorithm implementation
type lru struct {
	opts *lruOptions

	keys map[interface{}]*list.Element
	hot  *list.List
}

func newLRU(opts *lruOptions) *lru {
	return &lru{
		opts: opts,
		keys: make(map[interface{}]*list.Element, opts.capacity),
		hot:  list.New(),
	}
}
func (l *lru) ClearExpired() {
	var (
		ele  *list.Element
		v    *cacheValue
		hot  = l.hot
		keys = l.keys
	)
	for {
		ele = hot.Front()
		if ele == nil {
			break
		}
		v = ele.Value.(*cacheValue)
		if v.IsDeleted() {
			l.hot.Remove(ele)
			delete(keys, v.Key)
		} else {
			break
		}
	}
}

// Add the value to the cache, only when the key does not exist
func (l *lru) Add(key, value interface{}) (newkey bool) {
	ele, exists := l.keys[key]
	if exists {
		v := ele.Value.(*cacheValue)
		if v.IsDeleted() {
			newkey = true
			v.Value = value
			l.moveHot(ele)
		}
	} else {
		newkey = true
		l.add(key, value)
	}
	return
}
func (l *lru) add(key, value interface{}) {
	// capacity limit reached, pop front
	for l.hot.Len() >= l.opts.capacity {
		ele := l.hot.Front()
		v := ele.Value.(*cacheValue)
		delete(l.keys, v.Key)
		l.hot.Remove(ele)
	}
	// new value
	v := &cacheValue{
		Key:   key,
		Value: value,
	}
	if l.opts.expiry > 0 {
		v.Deadline = time.Now().Add(l.opts.expiry)
	}
	l.keys[key] = l.hot.PushBack(v)
	return
}
func (l *lru) moveHot(ele *list.Element) {
	v := ele.Value.(*cacheValue)
	l.hot.Remove(ele)
	if l.opts.expiry > 0 {
		v.Deadline = time.Now().Add(l.opts.expiry)
	}
	l.keys[v.Key] = l.hot.PushBack(v)
}

func (l *lru) Put(key, value interface{}) (newkey bool) {
	ele, exists := l.keys[key]
	if exists {
		// put
		v := ele.Value.(*cacheValue)
		if v.IsDeleted() {
			newkey = true
		}
		v.Value = value
		// move hot
		l.moveHot(ele)
	} else {
		newkey = true
		l.add(key, value)
	}
	return
}

// Get return cache value
func (l *lru) Get(key interface{}) (value interface{}, exists bool) {
	ele, exists := l.keys[key]
	if !exists {
		return
	}
	v := ele.Value.(*cacheValue)
	if v.IsDeleted() {
		delete(l.keys, key)
		l.hot.Remove(ele)
		exists = false
		return
	}
	value = v.Value

	// move hot
	l.moveHot(ele)
	return
}
func (l *lru) Delete(key ...interface{}) (changed int) {
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

func (l *lru) Len() int {
	return l.hot.Len()
}

func (l *lru) Clear() (e error) {
	l.hot.Init()
	for k := range l.keys {
		delete(l.keys, k)
	}
	return
}
