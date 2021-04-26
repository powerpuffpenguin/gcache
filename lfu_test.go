package gcache_test

import (
	"testing"
	"time"

	"github.com/powerpuffpenguin/gcache"
	"github.com/stretchr/testify/assert"
)

func TestLFU(t *testing.T) {
	var l gcache.Cache
	l = gcache.NewLFU(
		gcache.WithLFUCapacity(3),
	)
	// 0 2
	// 1 2
	// 2 2
	for i := 0; i < 3; i++ {
		l.Put(i, i)
		v, exists := l.Get(i)
		assert.True(t, exists)
		assert.Equal(t, v, i)
	}
	// 0 2+2
	// 1 2+1
	// 2 2
	i := 0
	v, exists := l.Get(i)
	assert.True(t, exists)
	assert.Equal(t, v, i)
	i = 0
	v, exists = l.Get(i)
	assert.True(t, exists)
	assert.Equal(t, v, i)
	i = 1
	v, exists = l.Get(i)
	assert.True(t, exists)
	assert.Equal(t, v, i)

	// 0 4
	// 1 3
	// 3 1
	l.Put(3, 3)
	v, exists = l.Get(2)
	assert.False(t, exists)

	// 0 4
	// 1 3
	// 4 1
	l.Put(4, 4)
	v, exists = l.Get(3)
	assert.False(t, exists)

	// 0 4
	// 1 3
	// 5 6
	for i := 0; i < 3; i++ {
		l.Put(5, 5)
		v, exists := l.Get(5)
		assert.True(t, exists)
		assert.Equal(t, v, 5)
	}

	// 0 4
	// 5 6
	// 6 1
	l.Put(6, 6)
	v, exists = l.Get(1)
	assert.False(t, exists)

	// 0 4
	// 5 6
	// 6 2
	v, exists = l.Get(6)
	assert.True(t, exists)
	assert.Equal(t, v, 6)

	// expire
	count := 3
	duration := time.Millisecond * 10
	l = gcache.NewLFU(
		gcache.WithLFUCapacity(count),
		gcache.WithLFUExpiry(duration),
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
	l = gcache.NewLFU(
		gcache.WithLFUCapacity(count+1),
		gcache.WithLFUExpiry(duration),
		gcache.WithLFUClear(duration),
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

}
