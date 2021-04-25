package gcache

import (
	"runtime"
	"time"
)

type LRU struct {
	*wrapper
}

func NewLRU(opt ...LRUOption) (lru *LRU) {
	opts := defaultLRUOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	w := &wrapper{
		impl:   NewLowLRU(opts.capacity, opts.expiry),
		closed: make(chan struct{}),
	}
	lru = &LRU{
		wrapper: w,
	}
	if opts.expiry > 0 {
		ticker := time.NewTicker(opts.clear)
		w.ticker = ticker
		go w.clearExpired(ticker.C)
		runtime.SetFinalizer(lru, (*LRU).stop)
	}
	return
}
func (l *LRU) stop() {
	l.wrapper.stop()
}
