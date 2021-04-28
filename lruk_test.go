package gcache_test

import (
	"testing"
	"time"

	"github.com/powerpuffpenguin/gcache"
	"github.com/stretchr/testify/assert"
)

func TestLRU_K3(t *testing.T) {
	var (
		l gcache.Cache
		h gcache.LowCache
	)
	h = gcache.NewLowLRU(gcache.WithLowLRUCapacity(3))
	// history value
	l = gcache.NewLRUK(
		gcache.WithLRUK(3),
		gcache.WithLRUKHistoryOnlyKey(false),
		gcache.WithLRUKHistory(h),
		gcache.WithLRUKCapacity(3),
	)
	for i := 0; i < 6; i++ {
		added := l.Add(i, i)
		assert.True(t, added)

		if i < 3 {
			assert.Equal(t, i+1, l.Len())

			val, exists := l.Get(i)
			assert.True(t, exists)
			assert.Equal(t, val, i)
			assert.Equal(t, 1, h.Len())

			val, exists = l.Get(i)
			assert.True(t, exists)
			assert.Equal(t, val, i)
			assert.Equal(t, 0, h.Len())
		} else {
			assert.Equal(t, i-3+1, h.Len())
		}
	}

	for j := 0; j < 2; j++ {
		for i := 0; i < 6; i++ {
			val, exists := l.Get(i)
			assert.True(t, exists)
			assert.Equal(t, val, i)
		}
	}
	for i := 0; i < 6; i++ {
		val, exists := l.Get(i)
		if i < 3 {
			assert.False(t, exists)
		} else {
			assert.True(t, exists)
			assert.Equal(t, val, i)
		}
	}
}
func TestLRU_K2(t *testing.T) {
	var (
		l gcache.Cache
		h gcache.LowCache
	)
	h = gcache.NewLowLRU(gcache.WithLowLRUCapacity(3))
	lru := gcache.NewLowLRU(
		gcache.WithLowLRUCapacity(3),
	)
	l = gcache.NewLRUK(
		gcache.WithLRUK(2),
		gcache.WithLRUKHistoryOnlyKey(true),
		gcache.WithLRUKHistory(h),
		gcache.WithLRUKCapacity(3),
	)
	count := 1000
	for i := 0; i < count; i++ {
		added := l.Add(i, i)
		assert.False(t, added)
		size := l.Len()
		assert.Equal(t, 0, size)
		size = h.Len()
		if i < 2 {
			assert.Equal(t, size, i+1)
		} else {
			assert.Equal(t, size, 3)
		}
		for j := 0; j < size; j++ {
			key := i - size + 1 + j
			_, exists := lru.Get(key)
			assert.False(t, exists)
		}
	}

	for i := 0; i < 4; i++ {
		added := l.Add(i, i)
		assert.False(t, added)
	}

	for i := 1; i < 4; i++ {
		added := l.Add(i, i)
		assert.True(t, added)

		size := l.Len()
		for j := 0; j < size; j++ {
			key := i - size + 1 + j
			val, exists := l.Get(key)
			assert.True(t, exists)
			assert.Equal(t, key, val)
		}
	}
	l.Clear()
	assert.Equal(t, 0, l.Len())
	assert.Equal(t, 0, h.Len())
	// history value
	l = gcache.NewLRUK(
		gcache.WithLRUK(2),
		gcache.WithLRUKHistoryOnlyKey(false),
		gcache.WithLRUKHistory(h),
		gcache.WithLRUKCapacity(3),
	)
	for i := 0; i < 6; i++ {
		added := l.Add(i, i)
		assert.True(t, added)
		if i < 3 {
			assert.Equal(t, i+1, l.Len())

			val, exists := l.Get(i)
			assert.True(t, exists)
			assert.Equal(t, val, i)

			assert.Equal(t, 0, h.Len())
		} else {
			assert.Equal(t, i-3+1, h.Len())
		}
	}
	for i := 0; i < 6; i++ {
		val, exists := l.Get(i)
		assert.True(t, exists)
		assert.Equal(t, val, i)
	}

	for i := 0; i < 6; i++ {
		val, exists := l.Get(i)
		if i < 3 {
			assert.False(t, exists)
		} else {
			assert.True(t, exists)
			assert.Equal(t, val, i)
		}
	}
}

func TestLRU_K1(t *testing.T) {
	// hot
	var l gcache.Cache
	l = gcache.NewLRUK(
		gcache.WithLRUK(1),
		gcache.WithLRUKCapacity(3),
	)
	count := 1000
	for i := 0; i < count; i++ {
		added := l.Add(i, i)
		assert.True(t, added)
		size := l.Len()
		if i < 2 {
			assert.Equal(t, size, i+1)
		} else {
			assert.Equal(t, size, 3)
		}
		for j := 0; j < size; j++ {
			key := i - size + 1 + j
			val, exists := l.Get(key)
			assert.True(t, exists)
			assert.Equal(t, key, val)
		}
	}
	// delete
	count = 3 * 30
	l = gcache.NewLRUK(
		gcache.WithLRUK(1),
		gcache.WithLRUKCapacity(count),
	)
	for i := 0; i < count; i++ {
		added := l.Add(i, i)
		assert.True(t, added)
		size := l.Len()
		assert.Equal(t, size, i+1)
	}
	for i := 0; i < count; i++ {
		if i%3 == 0 {
			changed := l.Delete(i, i+1)
			assert.Equal(t, changed, 2)
		}
	}
	for i := 0; i < count; i++ {
		if i%3 == 0 {
			key := i
			v, exists := l.Get(key)
			assert.False(t, exists)
			assert.Nil(t, v)

			key = i + 1
			v, exists = l.Get(key)
			assert.False(t, exists)
			assert.Nil(t, v)

			key = i + 2
			v, exists = l.Get(key)
			assert.True(t, exists)
			assert.Equal(t, key, v)

			changed := l.Delete(i, i+1)
			assert.Equal(t, changed, 0)
		}
	}

	// expire
	count = 3
	duration := time.Millisecond * 10 * 5
	l = gcache.NewLRUK(
		gcache.WithLRUK(1),
		gcache.WithLRUKCapacity(count),
		gcache.WithLRUKExpiry(duration),
	)
	for i := 0; i < count; i++ {
		l.Put(i, i)
		size := l.Len()
		assert.Equal(t, size, i+1)
	}
	for i := 0; i < count; i++ {
		key := i
		val, exists := l.Get(key)
		assert.True(t, exists)
		assert.Equal(t, key, val)
	}
	time.Sleep(duration)
	size := l.Len()
	assert.Equal(t, size, count)
	for i := 0; i < count; i++ {
		key := i
		v, exists := l.Get(key)
		assert.False(t, exists)
		assert.Nil(t, v)
	}

	// clear timer
	l = gcache.NewLRUK(
		gcache.WithLRUK(1),
		gcache.WithLRUKCapacity(count+1),
		gcache.WithLRUKExpiry(duration),
		gcache.WithLRUKClear(duration),
	)
	for i := 0; i < count; i++ {
		l.Put(i, i)
		size := l.Len()
		assert.Equal(t, size, i+1)
	}
	time.Sleep(duration / 2)
	l.Put("ok", "value")
	size = l.Len()
	assert.Equal(t, size, count+1)
	time.Sleep(duration/2 + duration/3)

	size = l.Len()
	assert.Equal(t, size, 1)

	key := "ok"
	val, exists := l.Get(key)
	assert.True(t, exists)
	assert.Equal(t, val, "value")

	time.Sleep(duration / 2)
	size = l.Len()
	assert.Equal(t, size, 1)
	key = "ok"
	val, exists = l.Get(key)
	assert.True(t, exists)
	assert.Equal(t, val, "value")

	time.Sleep(duration * 2)
	size = l.Len()
	assert.Equal(t, size, 0)

	// batch
	l.BatchPut(1, "1", 2)
	size = l.Len()
	assert.Equal(t, size, 2)
	vals := l.BatchGet(1, 2, 3)
	assert.Equal(t, len(vals), 3)
	assert.True(t, vals[0].Exists)
	assert.Equal(t, vals[0].Value, "1")
	assert.True(t, vals[1].Exists)
	assert.Nil(t, vals[1].Value)
	assert.False(t, vals[2].Exists)
	assert.Nil(t, vals[2].Value)
}
