package gcache

import "time"

var defaultFIFOOptions = fifoOptions{
	expiry:   0,
	capacity: 1000,
	clear:    time.Minute * 10,
}

type fifoOptions struct {
	expiry   time.Duration
	capacity int
	clear    time.Duration
}
type FIFOOption interface {
	apply(*fifoOptions)
}
type funcFIFOOption struct {
	f func(*fifoOptions)
}

func (fdo *funcFIFOOption) apply(do *fifoOptions) {
	fdo.f(do)
}
func newFuncFIFOOption(f func(*fifoOptions)) *funcFIFOOption {
	return &funcFIFOOption{
		f: f,
	}
}

// WithFIFOExpiry if <=0, it will not expire due to time
func WithFIFOExpiry(expiry time.Duration) FIFOOption {
	return newFuncFIFOOption(func(o *fifoOptions) {
		o.expiry = expiry
	})
}

// WithFIFOCapacity set the maximum amount of data to be cached
func WithFIFOCapacity(capacity int) FIFOOption {
	return newFuncFIFOOption(func(o *fifoOptions) {
		if capacity < 1 {
			panic(`fifo capacity must > 0`)
		}
		o.capacity = capacity
	})
}

// WithFIFOClear timer clear expired cache, if <=0 not start timer.
func WithFIFOClear(duration time.Duration) FIFOOption {
	return newFuncFIFOOption(func(po *fifoOptions) {
		po.clear = duration
	})
}
