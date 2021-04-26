package gcache

// NewLowFIFO create a low-level lru, use NewFIFO unless you know exactly what you are doing.
func NewLowFIFO(opt ...LowFIFOOption) LowCache {
	opts := defaultLowFIFOOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	return newLRUFIFO(false,
		opts.capacity,
		opts.expiry,
	)
}
