package gcache_test

import (
	"testing"

	"github.com/powerpuffpenguin/gcache"
	"github.com/stretchr/testify/assert"
)

func TestFIFO(t *testing.T) {
	// fifo
	var l gcache.Cache
	l = gcache.NewFIFO(
		gcache.WithFIFOCapacity(3),
	)
	for i := 0; i < 3; i++ {
		l.Put(i, i)
		v, exists := l.Get(i)
		assert.True(t, exists)
		assert.Equal(t, v, i)
	}
	v, exists := l.Get(0)
	assert.True(t, exists)
	assert.Equal(t, v, 0)
	l.Put(3, 3)
	for i := 0; i < 4; i++ {
		v, exists := l.Get(i)
		if i == 0 {
			assert.False(t, exists)
		} else {
			assert.True(t, exists)
			assert.Equal(t, v, i)
		}
	}

	// lru
	l = gcache.NewLRU(
		gcache.WithLRUCapacity(3),
	)
	for i := 0; i < 3; i++ {
		l.Put(i, i)
		v, exists := l.Get(i)
		assert.True(t, exists)
		assert.Equal(t, v, i)
	}
	v, exists = l.Get(0)
	assert.True(t, exists)
	assert.Equal(t, v, 0)
	l.Put(3, 3)
	for i := 0; i < 4; i++ {
		v, exists := l.Get(i)
		if i == 1 {
			assert.False(t, exists)
		} else {
			assert.True(t, exists)
			assert.Equal(t, v, i)
		}
	}
}
