package gcache

import (
	"runtime"
	"time"
)

type FIFO struct {
	*wrapper
}

func NewFIFO(opt ...FIFOOption) (fifo *FIFO) {
	opts := defaultFIFOOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	w := &wrapper{
		impl: NewLowFIFO(
			WithLowFIFOCapacity(opts.capacity),
			WithLowFIFOExpiry(opts.expiry),
		),
		closed: make(chan struct{}),
	}
	fifo = &FIFO{
		wrapper: w,
	}
	if opts.expiry > 0 {
		ticker := time.NewTicker(opts.clear)
		w.ticker = ticker
		go w.clearExpired(ticker.C)
		runtime.SetFinalizer(fifo, (*FIFO).stop)
	}
	return
}
func (f *FIFO) stop() {
	f.wrapper.stop()
}
