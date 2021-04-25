package gcache

import "time"

var defaultLowLRUOptions = lowLRUOptions{
	expiry:   0,
	capacity: 1000,
}

type lowLRUOptions struct {
	expiry   time.Duration
	capacity int
}
type LowLRUOption interface {
	apply(*lowLRUOptions)
}
type funcLowLRUOption struct {
	f func(*lowLRUOptions)
}

func (fdo *funcLowLRUOption) apply(do *lowLRUOptions) {
	fdo.f(do)
}
func newFuncLowLRUOption(f func(*lowLRUOptions)) *funcLowLRUOption {
	return &funcLowLRUOption{
		f: f,
	}
}

// WithLowLRUExpiry if <=0, it will not expire due to time
func WithLowLRUExpiry(expiry time.Duration) LowLRUOption {
	return newFuncLowLRUOption(func(o *lowLRUOptions) {
		o.expiry = expiry
	})
}

// WithLRUCapacity set the maximum amount of data to be cached
func WithLowLRUCapacity(capacity int) LowLRUOption {
	return newFuncLowLRUOption(func(o *lowLRUOptions) {
		if capacity < 1 {
			panic(`lru capacity must > 0`)
		}
		o.capacity = capacity
	})
}
