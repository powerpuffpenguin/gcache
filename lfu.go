package gcache

import (
	"runtime"
	"time"
)

type LFU struct {
	*wrapper
}

func NewLFU(opt ...LFUOption) (lfu *LFU) {
	opts := defaultLFUOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	w := &wrapper{
		impl: NewLowLFU(
			WithLowLFUCapacity(opts.capacity),
			WithLowLFUExpiry(opts.expiry),
		),
		closed: make(chan struct{}),
	}
	lfu = &LFU{
		wrapper: w,
	}
	if opts.expiry > 0 {
		ticker := time.NewTicker(opts.clear)
		w.ticker = ticker
		go w.clearExpired(ticker.C)
		runtime.SetFinalizer(lfu, (*LFU).stop)
	}
	return
}
func (l *LFU) stop() {
	l.wrapper.stop()
}
