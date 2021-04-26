package gcache

import "time"

var defaultLowLFUOptions = lowLFUOptions{
	expiry:   0,
	capacity: 1000,
}

type lowLFUOptions struct {
	expiry   time.Duration
	capacity int
}
type LowLFUOption interface {
	apply(*lowLFUOptions)
}
type funcLowLFUOption struct {
	f func(*lowLFUOptions)
}

func (fdo *funcLowLFUOption) apply(do *lowLFUOptions) {
	fdo.f(do)
}
func newFuncLowLFUOption(f func(*lowLFUOptions)) *funcLowLFUOption {
	return &funcLowLFUOption{
		f: f,
	}
}

// WithLowLFUExpiry if <=0, it will not expire due to time
func WithLowLFUExpiry(expiry time.Duration) LowLFUOption {
	return newFuncLowLFUOption(func(o *lowLFUOptions) {
		o.expiry = expiry
	})
}

// WithLowLFUCapacity set the maximum amount of data to be cached
func WithLowLFUCapacity(capacity int) LowLFUOption {
	return newFuncLowLFUOption(func(o *lowLFUOptions) {
		if capacity < 1 {
			panic(`lfu capacity must > 0`)
		}
		o.capacity = capacity
	})
}
