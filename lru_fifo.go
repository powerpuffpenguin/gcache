package gcache

import (
	"container/list"
	"time"
)

type lrufifo struct {
	keys     map[interface{}]*list.Element
	hot      *list.List
	expiry   time.Duration
	capacity int
	lru      bool
}

func newLRUFIFO(lru bool, capacity int, expiry time.Duration) *lrufifo {
	return &lrufifo{
		keys:     make(map[interface{}]*list.Element, capacity),
		hot:      list.New(),
		expiry:   expiry,
		capacity: capacity,
		lru:      lru,
	}
}

func (l *lrufifo) ClearExpired() {
	if l.expiry > 0 {
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
				hot.Remove(ele)
				delete(keys, v.GetKey())
			} else {
				break
			}
		}
	}
}

// Add the value to the cache, only when the key does not exist
func (l *lrufifo) Add(key, value interface{}) (added bool) {
	ele, exists := l.keys[key]
	if exists {
		v := ele.Value.(cacheValue)
		if v.IsDeleted() {
			added = true
			v.SetValue(value)
			l.moveHot(ele)
			l.ClearExpired()
		}
	} else {
		added = true
		l.add(key, value)
	}
	return
}

func (l *lrufifo) add(key, value interface{}) (delkey, delval interface{}, deleted bool) {
	// capacity limit reached, pop front
	if l.hot.Len() >= l.capacity {
		deleted = true
		ele := l.hot.Front()
		v := ele.Value.(cacheValue)
		delkey = v.GetKey()
		delval = v.GetValue()
		delete(l.keys, delkey)
		l.hot.Remove(ele)
	}
	// new value
	v := newValue(key, value, l.expiry)
	l.keys[key] = l.hot.PushBack(v)
	return
}

func (l *lrufifo) moveHot(ele *list.Element) {
	v := ele.Value.(cacheValue)
	if l.expiry > 0 {
		v.SetDeadline(time.Now().Add(l.expiry))
	}
	l.hot.MoveToBack(ele)
}

func (l *lrufifo) Put(key, value interface{}) (delkey, delval interface{}, deleted bool) {
	ele, exists := l.keys[key]
	if exists {
		// put
		v := ele.Value.(cacheValue)
		if v.IsDeleted() {
			v.SetValue(value)
			// move hot
			l.moveHot(ele)

			l.ClearExpired()
		} else {
			deleted = true
			delkey = key
			delval = v.GetValue()

			v.SetValue(value)
			// move hot
			l.moveHot(ele)
		}

	} else {
		delkey, delval, deleted = l.add(key, value)
	}
	return
}

// Get return cache value
func (l *lrufifo) Get(key interface{}) (value interface{}, exists bool) {
	ele, exists := l.keys[key]
	if !exists {
		return
	}
	v := ele.Value.(cacheValue)
	if v.IsDeleted() {
		delete(l.keys, key)
		l.hot.Remove(ele)
		exists = false
		l.ClearExpired()
		return
	}
	value = v.GetValue()

	// fifo not need move hot
	if l.lru {
		l.moveHot(ele)
	}
	return
}

func (l *lrufifo) Delete(key ...interface{}) (changed int) {
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

func (l *lrufifo) Len() int {
	return l.hot.Len()
}

func (l *lrufifo) Clear() {
	l.hot.Init()
	for k := range l.keys {
		delete(l.keys, k)
	}
	return
}
