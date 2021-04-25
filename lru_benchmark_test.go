package gcache_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/powerpuffpenguin/gcache"
)

func BenchmarkLruNotFoundAdd(b *testing.B) {
	l := gcache.NewLRU(
		gcache.WithLRUCapacity(b.N),
	)

	var finish sync.WaitGroup
	var added int32
	var idle int32

	fn := func(id int) {
		for i := 0; i < b.N; i++ {
			add := l.Add(i, i+id)
			if add {
				atomic.AddInt32(&added, 1)
			} else {
				atomic.AddInt32(&idle, 1)
			}
			time.Sleep(0)
		}
		finish.Done()
	}

	finish.Add(10)
	go fn(0x0000)
	go fn(0x1100)
	go fn(0x2200)
	go fn(0x3300)
	go fn(0x4400)
	go fn(0x5500)
	go fn(0x6600)
	go fn(0x7700)
	go fn(0x8800)
	go fn(0x9900)
	finish.Wait()
}
