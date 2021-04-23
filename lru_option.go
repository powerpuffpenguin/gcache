package gcache

import "time"

var defaultLRUOptions = lruOptions{
	expiry:   0,
	capacity: 1000,
	clear:    time.Minute * 10,
}

type lruOptions struct {
	expiry   time.Duration
	capacity int
	clear    time.Duration
}
type LRUOption interface {
	apply(*lruOptions)
}
type funcLRUOption struct {
	f func(*lruOptions)
}

func (fdo *funcLRUOption) apply(do *lruOptions) {
	fdo.f(do)
}
func newFuncLRUOption(f func(*lruOptions)) *funcLRUOption {
	return &funcLRUOption{
		f: f,
	}
}

// WithLRUExpiry if <=0, it will not expire due to time
func WithLRUExpiry(expiry time.Duration) LRUOption {
	return newFuncLRUOption(func(o *lruOptions) {
		o.expiry = expiry
	})
}

// WithLRUCapacity set the maximum amount of data to be cached
func WithLRUCapacity(capacity int) LRUOption {
	return newFuncLRUOption(func(o *lruOptions) {
		if capacity < 1 {
			panic(`lru capacity must > 0`)
		}
		o.capacity = capacity
	})
}

// WithLRUClear timer clear expired cache, if <=0 not start timer.
func WithLRUClear(duration time.Duration) LRUOption {
	return newFuncLRUOption(func(po *lruOptions) {
		po.clear = duration
	})
}
