package gcache

import "time"

var defaultLRUKOptions = lrukOptions{
	historyOnlyKey: true,
	expiry:         0,
	capacity:       1000,
	clear:          time.Minute * 10,
	k:              2,
}

type lrukOptions struct {
	lru, history   LowCache
	historyOnlyKey bool
	expiry         time.Duration
	capacity       int
	clear          time.Duration
	k              int
}
type LRUKOption interface {
	apply(*lrukOptions)
}
type funcLRUKOption struct {
	f func(*lrukOptions)
}

func (fdo *funcLRUKOption) apply(do *lrukOptions) {
	fdo.f(do)
}
func newFuncLRUKOption(f func(*lrukOptions)) *funcLRUKOption {
	return &funcLRUKOption{
		f: f,
	}
}

// WithLRUKExpiry if <=0, it will not expire due to time
func WithLRUKExpiry(expiry time.Duration) LRUKOption {
	return newFuncLRUKOption(func(o *lrukOptions) {
		o.expiry = expiry
	})
}

// WithLRUKCapacity set the maximum amount of data to be cached
func WithLRUKCapacity(capacity int) LRUKOption {
	return newFuncLRUKOption(func(o *lrukOptions) {
		if capacity < 1 {
			panic(`lru capacity must > 0`)
		}
		o.capacity = capacity
	})
}

// WithLRUKClear timer clear expired cache, if <=0 not start timer.
func WithLRUKClear(duration time.Duration) LRUKOption {
	return newFuncLRUKOption(func(po *lrukOptions) {
		po.clear = duration
	})
}

// WithLRUK set lru-k ,if k == 1 use lru, if k >1 use lru-k, if < 1
func WithLRUK(k int) LRUKOption {
	return newFuncLRUKOption(func(po *lrukOptions) {
		if k < 1 {
			panic("lru-k k must > 0")
		}
		po.k = k
	})
}

// WithLRUKLRU if lru is nil auto create.
func WithLRUKLRU(lru LowCache) LRUKOption {
	return newFuncLRUKOption(func(po *lrukOptions) {
		po.lru = lru
	})
}

// WithLRUKHistory if history is nil use lru-1, default is nil.
func WithLRUKHistory(history LowCache) LRUKOption {
	return newFuncLRUKOption(func(po *lrukOptions) {
		po.history = history
	})
}

// WithLRUKHistoryOnlyKey if ture history only save key, if false history will save key and value
func WithLRUKHistoryOnlyKey(onlyKey bool) LRUKOption {
	return newFuncLRUKOption(func(po *lrukOptions) {
		po.historyOnlyKey = onlyKey
	})
}
