package gcache

import "time"

var defaultLFUOptions = lfuOptions{
	expiry:   0,
	capacity: 1000,
	clear:    time.Minute * 10,
}

type lfuOptions struct {
	expiry   time.Duration
	capacity int
	clear    time.Duration
}
type LFUOption interface {
	apply(*lfuOptions)
}
type funcLFUOption struct {
	f func(*lfuOptions)
}

func (fdo *funcLFUOption) apply(do *lfuOptions) {
	fdo.f(do)
}
func newFuncLFUOption(f func(*lfuOptions)) *funcLFUOption {
	return &funcLFUOption{
		f: f,
	}
}

// WithLFUExpiry if <=0, it will not expire due to time
func WithLFUExpiry(expiry time.Duration) LFUOption {
	return newFuncLFUOption(func(o *lfuOptions) {
		o.expiry = expiry
	})
}

// WithLFUCapacity set the maximum amount of data to be cached
func WithLFUCapacity(capacity int) LFUOption {
	return newFuncLFUOption(func(o *lfuOptions) {
		if capacity < 1 {
			panic(`lfu capacity must > 0`)
		}
		o.capacity = capacity
	})
}

// WithLFUClear timer clear expired cache, if <=0 not start timer.
func WithLFUClear(duration time.Duration) LFUOption {
	return newFuncLFUOption(func(po *lfuOptions) {
		po.clear = duration
	})
}
