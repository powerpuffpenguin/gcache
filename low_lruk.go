package gcache

type kValue struct {
	Count int
	Key   interface{}
	Value interface{}
}

// A low-level implementation of lruk, use LRUK unless you know exactly what you are doing.
type LowLRUK struct {
	opts         lowLRUKOptions
	history, lru LowCache
}

// NewLowLRUK create a low-level lru, use NewLRUK unless you know exactly what you are doing.
func NewLowLRUK(history, lru LowCache, opt ...LowLRUKOption) *LowLRUK {
	opts := defaultLowLRUKOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	if opts.k < 2 {
		history = nil
	}
	return &LowLRUK{
		opts:    opts,
		lru:     lru,
		history: history,
	}
}

// Clear Expired cache
func (l *LowLRUK) ClearExpired() {
	l.lru.ClearExpired()
	if l.history != nil {
		l.history.ClearExpired()
	}
}

// Add the value to the cache, only when the key does not exist
func (l *LowLRUK) Add(key, value interface{}) (added bool) {
	_, exists := l.lru.Get(key)
	if exists {
		return
	} else if l.history == nil {
		added = l.lru.Add(key, value)
		return
	}

	v, exists := l.history.Get(key)
	if exists {
		kv := v.(kValue)
		kv.Count++
		if kv.Count >= l.opts.k {
			l.history.Delete(key)
			added = l.lru.Add(key, value)
		} else {
			l.history.Put(key, kv)
		}
	} else {
		kv := kValue{
			Count: 1,
			Key:   key,
		}
		if !l.opts.historyOnlyKey {
			kv.Value = value
			added = true
		}
		l.history.Put(key, kv)
	}
	return
}

// Put key value to cache
func (l *LowLRUK) Put(key, value interface{}) (delkey, delval interface{}, deleted bool) {
	_, exists := l.lru.Get(key)
	if exists {
		delkey, delval, deleted = l.lru.Put(key, value)
		return
	} else if l.history == nil {
		delkey, delval, deleted = l.lru.Put(key, value)
		return
	}

	v, exists := l.history.Get(key)
	if exists {
		kv := v.(kValue)
		kv.Count++
		if kv.Count >= l.opts.k {
			l.history.Delete(key)
			delkey, delval, deleted = l.lru.Put(key, value)
		} else {
			l.history.Put(key, kv)
		}
	} else {
		kv := kValue{
			Count: 1,
			Key:   key,
		}
		if l.opts.historyOnlyKey {
			l.history.Put(key, kv)
		} else {
			kv.Value = value
			delkey, delval, deleted = l.history.Put(key, kv)
		}
	}
	return
}

// Get return cache value
func (l *LowLRUK) Get(key interface{}) (value interface{}, exists bool) {
	value, exists = l.lru.Get(key)
	if exists || l.history == nil {
		return
	}
	v, ok := l.history.Get(key)
	if ok {
		kv := v.(kValue)
		value = kv.Value

		if l.opts.historyOnlyKey {
			// mov to hot
			if kv.Count < l.opts.k-1 {
				kv.Count++
			}
			l.history.Put(key, kv)
		} else {
			exists = true
			kv.Count++
			if kv.Count >= l.opts.k {
				l.history.Delete(key)
				l.lru.Put(key, kv.Value)
			} else {
				l.history.Put(key, kv)
			}
		}
	} else if l.opts.historyOnlyKey {
		l.history.Put(key, kValue{
			Count: 1,
			Key:   key,
		})
	}
	return
}

// Delete key from cache
func (l *LowLRUK) Delete(key ...interface{}) (changed int) {
	changed = l.lru.Delete(key...)
	if l.history != nil {
		if l.opts.historyOnlyKey {
			l.history.Delete(key...)
		} else {
			changed += l.history.Delete(key...)
		}
	}
	return
}

// Len returns the number of cached data
func (l *LowLRUK) Len() int {
	count := l.lru.Len()
	if l.history != nil && !l.opts.historyOnlyKey {
		count += l.history.Len()
	}
	return count
}

// Clear all cached data
func (l *LowLRUK) Clear() {
	l.lru.Clear()
	if l.history != nil {
		l.history.Clear()
	}
}
