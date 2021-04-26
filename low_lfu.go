package gcache

// NewLowLFU create a low-level lfu, use NewLFU unless you know exactly what you are doing.
func NewLowLFU(opt ...LowLFUOption) LowCache {
	opts := defaultLowLFUOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	if opts.expiry > 0 {
		// return newLowLFUEx(opts.capacity, opts.expiry)
	}
	return newLowLFU(opts.capacity)
}

type lowLFU struct {
	keys     map[interface{}]lfuValue
	hot      *lfuHeap
	capacity int
}

func newLowLFU(capacity int) *lowLFU {
	return &lowLFU{
		keys:     make(map[interface{}]lfuValue, capacity),
		hot:      newLFUHeap(capacity),
		capacity: capacity,
	}
}
func (l *lowLFU) ClearExpired() {
}

// Add the value to the cache, only when the key does not exist
func (l *lowLFU) Add(key, value interface{}) (added bool) {
	_, exists := l.keys[key]
	if !exists {
		added = true
		l.add(key, value)
	}
	return
}

func (l *lowLFU) add(key, value interface{}) (delkey, delval interface{}, deleted bool) {
	// capacity limit reached, pop
	if l.hot.Len() >= l.capacity {
		deleted = true
		v := l.hot.Remove(0)
		delkey = v.GetKey()
		delval = v.GetValue()
		delete(l.keys, delkey)
	}
	// new value
	v := newLFUValue(key, value, 0)
	l.keys[key] = v
	l.hot.Push(v)
	return
}
func (l *lowLFU) moveHot(v lfuValue) {
	v.Increment()
	l.hot.Fix(v.GetIndex())
}
func (l *lowLFU) Put(key, value interface{}) (delkey, delval interface{}, deleted bool) {
	v, exists := l.keys[key]
	if exists {
		deleted = true
		delkey = key
		delval = v.GetValue()

		// put
		v.SetKey(value)
		// move hot
		l.moveHot(v)
	} else {
		delkey, delval, deleted = l.add(key, value)
	}
	return
}

// Get return cache value
func (l *lowLFU) Get(key interface{}) (value interface{}, exists bool) {
	v, exists := l.keys[key]
	if !exists {
		return
	}
	value = v.GetValue()

	// move hot
	l.moveHot(v)
	return
}
func (l *lowLFU) Delete(key ...interface{}) (changed int) {
	var (
		v      lfuValue
		exists bool
	)
	for _, k := range key {
		v, exists = l.keys[k]
		if exists {
			changed++
			delete(l.keys, k)
			l.hot.Remove(v.GetIndex())
		}
	}
	return
}

func (l *lowLFU) Len() int {
	return l.hot.Len()
}

func (l *lowLFU) Clear() {
	l.hot.Clear()
	for k := range l.keys {
		delete(l.keys, k)
	}
	return
}
