package gcache

import (
	"container/list"
	"time"
)

type lowLFUEx struct {
	keys     map[interface{}]*list.Element
	hot      *lfuHeap
	list     *list.List
	capacity int
	expiry   time.Duration
}

func newLowLFUEx(capacity int, expiry time.Duration) *lowLFUEx {
	return &lowLFUEx{
		keys:     make(map[interface{}]*list.Element, capacity),
		hot:      newLFUHeap(capacity),
		list:     list.New(),
		capacity: capacity,
		expiry:   expiry,
	}
}
func (l *lowLFUEx) ClearExpired() {
	if l.expiry > 0 {
		var (
			ele  *list.Element
			v    lfuValue
			list = l.list
			hot  = l.hot
			keys = l.keys
		)
		for {
			ele = list.Front()
			if ele == nil {
				break
			}
			v = ele.Value.(lfuValue)
			if v.IsDeleted() {
				list.Remove(ele)
				hot.Remove(v.GetIndex())
				delete(keys, v.GetKey())
			} else {
				break
			}
		}
	}
}

// Add the value to the cache, only when the key does not exist
func (l *lowLFUEx) Add(key, value interface{}) (added bool) {
	ele, exists := l.keys[key]
	if exists {
		v := ele.Value.(lfuValue)
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
func (l *lowLFUEx) add(key, value interface{}) (delkey, delval interface{}, deleted bool) {
	// capacity limit reached, pop front
	if l.hot.Len() >= l.capacity {
		deleted = true
		v := l.hot.heap[0]
		delkey = v.GetKey()
		delval = v.GetValue()
		delete(l.keys, delkey)
		l.list.Remove(l.keys[delkey])
		l.hot.Remove(0)
	}
	// new value
	v := newLFUValue(key, value, l.expiry)
	l.keys[key] = l.list.PushBack(v)
	l.hot.Push(v)
	return
}
func (l *lowLFUEx) moveHot(ele *list.Element) {
	v := ele.Value.(lfuValue)
	v.SetDeadline(time.Now().Add(l.expiry))
	l.list.MoveToBack(ele)

	v.Increment()
	l.hot.Fix(v.GetIndex())
}
func (l *lowLFUEx) Put(key, value interface{}) (delkey, delval interface{}, deleted bool) {
	ele, exists := l.keys[key]
	if exists {
		// put
		v := ele.Value.(lfuValue)
		if v.IsDeleted() {
			v.SetKey(value)
			// move hot
			l.moveHot(ele)

			l.ClearExpired()
		} else {
			deleted = true
			delkey = key
			delval = v.GetValue()

			v.SetKey(value)
			// move hot
			l.moveHot(ele)
		}

	} else {
		delkey, delval, deleted = l.add(key, value)
	}
	return
}

// Get return cache value
func (l *lowLFUEx) Get(key interface{}) (value interface{}, exists bool) {
	ele, exists := l.keys[key]
	if !exists {
		return
	}
	v := ele.Value.(lfuValue)
	if v.IsDeleted() {
		delete(l.keys, key)
		l.list.Remove(ele)
		l.hot.Remove(v.GetIndex())
		exists = false
		l.ClearExpired()
		return
	}
	value = v.GetValue()

	// move hot
	l.moveHot(ele)
	return
}
func (l *lowLFUEx) Delete(key ...interface{}) (changed int) {
	var (
		ele    *list.Element
		exists bool
		v      lfuValue
	)
	for _, k := range key {
		ele, exists = l.keys[k]
		if exists {
			v = ele.Value.(lfuValue)
			changed++
			delete(l.keys, k)
			l.list.Remove(ele)
			l.hot.Remove(v.GetIndex())
		}
	}
	return
}

func (l *lowLFUEx) Len() int {
	return l.hot.Len()
}

func (l *lowLFUEx) Clear() {
	l.list.Init()
	l.hot.Clear()
	for k := range l.keys {
		delete(l.keys, k)
	}
	return
}
