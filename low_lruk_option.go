package gcache

var defaultLowLRUKOptions = lowLRUKOptions{
	k:              2,
	historyOnlyKey: true,
}

type lowLRUKOptions struct {
	k              int
	historyOnlyKey bool
}
type LowLRUKOption interface {
	apply(*lowLRUKOptions)
}
type funcLowLRUKOption struct {
	f func(*lowLRUKOptions)
}

func (fdo *funcLowLRUKOption) apply(do *lowLRUKOptions) {
	fdo.f(do)
}
func newFuncLowLRUKOption(f func(*lowLRUKOptions)) *funcLowLRUKOption {
	return &funcLowLRUKOption{
		f: f,
	}
}

// WithLowLRUK set lru-k ,if k == 1 use lru, if k >1 use lru-k, if < 1
func WithLowLRUK(k int) LowLRUKOption {
	return newFuncLowLRUKOption(func(po *lowLRUKOptions) {
		if k < 1 {
			panic("lru-k k must > 0")
		}
		po.k = k
	})
}

// WithLowLRUKHistoryOnlyKey if ture history only save key, if false history will save key and value
func WithLowLRUKHistoryOnlyKey(onlyKey bool) LowLRUKOption {
	return newFuncLowLRUKOption(func(po *lowLRUKOptions) {
		po.historyOnlyKey = onlyKey
	})
}
