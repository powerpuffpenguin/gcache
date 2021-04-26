package gcache

import "time"

var defaultLowFIFOOptions = lowFIFOOptions{
	expiry:   0,
	capacity: 1000,
}

type lowFIFOOptions struct {
	expiry   time.Duration
	capacity int
}
type LowFIFOOption interface {
	apply(*lowFIFOOptions)
}
type funcLowFIFOOption struct {
	f func(*lowFIFOOptions)
}

func (fdo *funcLowFIFOOption) apply(do *lowFIFOOptions) {
	fdo.f(do)
}
func newFuncLowFIFOOption(f func(*lowFIFOOptions)) *funcLowFIFOOption {
	return &funcLowFIFOOption{
		f: f,
	}
}

// WithLowFIFOExpiry if <=0, it will not expire due to time
func WithLowFIFOExpiry(expiry time.Duration) LowFIFOOption {
	return newFuncLowFIFOOption(func(o *lowFIFOOptions) {
		o.expiry = expiry
	})
}

// WithLowFIFOCapacity set the maximum amount of data to be cached
func WithLowFIFOCapacity(capacity int) LowFIFOOption {
	return newFuncLowFIFOOption(func(o *lowFIFOOptions) {
		if capacity < 1 {
			panic(`fifo capacity must > 0`)
		}
		o.capacity = capacity
	})
}
