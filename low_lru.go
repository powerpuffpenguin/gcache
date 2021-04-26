package gcache

// NewLowLRU create a low-level lru, use NewLRU unless you know exactly what you are doing.
func NewLowLRU(opt ...LowLRUOption) LowCache {
	opts := defaultLowLRUOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	return newLRUFIFO(true,
		opts.capacity,
		opts.expiry,
	)
}
