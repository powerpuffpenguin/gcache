package gcache

import (
	"runtime"
	"time"
)

type LRUK struct {
	*wrapper
}

func NewLRUK(opt ...LRUKOption) (lruk *LRUK) {
	opts := defaultLRUKOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	// create default lru
	if opts.lru == nil {
		opts.lru = NewLowLRU(
			WithLowLRUCapacity(opts.capacity),
			WithLowLRUExpiry(opts.expiry),
		)
	}
	// create default history
	if opts.k > 1 && opts.history == nil {
		capacity := opts.capacity
		if opts.historyOnlyKey {
			capacity *= 10
		}
		opts.history = NewLowLRU(
			WithLowLRUCapacity(capacity),
			WithLowLRUExpiry(opts.expiry),
		)
	}

	w := &wrapper{
		impl: NewLowLRUK(
			opts.history, opts.lru,
			WithLowLRUK(opts.k),
			WithLowLRUKHistoryOnlyKey(opts.historyOnlyKey),
		),
		closed: make(chan struct{}),
	}
	lruk = &LRUK{
		wrapper: w,
	}
	if opts.expiry > 0 {
		ticker := time.NewTicker(opts.clear)
		lruk.ticker = ticker
		go w.clearExpired(ticker.C)
		runtime.SetFinalizer(lruk, (*LRUK).stop)
	}
	return
}
func (l *LRUK) stop() {
	l.wrapper.stop()
}
