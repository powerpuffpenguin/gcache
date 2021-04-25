package gcache_test

import (
	"testing"
	"time"

	"github.com/powerpuffpenguin/gcache"
	"github.com/stretchr/testify/assert"
)

func TestLRU(t *testing.T) {
	// hot
	l := gcache.NewLRU(
		gcache.WithLRUCapacity(3),
	)
	count := 1000
	for i := 0; i < count; i++ {
		added, e := l.Add(i, i)
		assert.Nil(t, e)
		assert.True(t, added)
		size, e := l.Len()
		assert.Nil(t, e)
		if i < 2 {
			assert.Equal(t, size, i+1)
		} else {
			assert.Equal(t, size, 3)
		}
		for j := 0; j < size; j++ {
			key := i - size + 1 + j
			val, e := l.Get(key)
			assert.Nil(t, e)
			assert.Equal(t, key, val)
		}
	}

	// delete
	count = 3 * 30
	l = gcache.NewLRU(
		gcache.WithLRUCapacity(count),
	)
	for i := 0; i < count; i++ {
		added, e := l.Add(i, i)
		assert.Nil(t, e)
		assert.True(t, added)
		size, e := l.Len()
		assert.Nil(t, e)
		assert.Equal(t, size, i+1)
	}
	for i := 0; i < count; i++ {
		if i%3 == 0 {
			changed, e := l.Delete(i, i+1)
			assert.Nil(t, e)
			assert.Equal(t, changed, 2)
		}
	}
	for i := 0; i < count; i++ {
		if i%3 == 0 {
			key := i
			v, e := l.Get(key)
			assert.ErrorIs(t, e, gcache.ErrNotExists)
			assert.Nil(t, v)

			key = i + 1
			v, e = l.Get(key)
			assert.ErrorIs(t, e, gcache.ErrNotExists)
			assert.Nil(t, v)

			key = i + 2
			v, e = l.Get(key)
			assert.Nil(t, e)
			assert.Equal(t, key, v)

			changed, e := l.Delete(i, i+1)
			assert.Nil(t, e)
			assert.Equal(t, changed, 0)
		}
	}

	// expire
	count = 3
	duration := time.Millisecond * 10

	l = gcache.NewLRU(
		gcache.WithLRUCapacity(count),
		gcache.WithLRUExpiry(duration),
	)
	for i := 0; i < count; i++ {
		added, e := l.Put(i, i)
		assert.Nil(t, e)
		assert.True(t, added)
		size, e := l.Len()
		assert.Nil(t, e)
		assert.Equal(t, size, i+1)
	}
	for i := 0; i < count; i++ {
		key := i
		val, e := l.Get(key)
		assert.Nil(t, e)
		assert.Equal(t, key, val)
	}
	time.Sleep(duration)
	size, e := l.Len()
	assert.Nil(t, e)
	assert.Equal(t, size, count)
	for i := 0; i < count; i++ {
		key := i
		v, e := l.Get(key)
		assert.ErrorIs(t, e, gcache.ErrNotExists)
		assert.Nil(t, v)
	}

	// clear timer
	l = gcache.NewLRU(
		gcache.WithLRUCapacity(count+1),
		gcache.WithLRUExpiry(duration),
		gcache.WithLRUClear(duration),
	)
	for i := 0; i < count; i++ {
		added, e := l.Put(i, i)
		assert.Nil(t, e)
		assert.True(t, added)
		size, e := l.Len()
		assert.Nil(t, e)
		assert.Equal(t, size, i+1)
	}
	time.Sleep(duration / 2)
	added, e := l.Put("ok", "value")
	assert.Nil(t, e)
	assert.True(t, added)
	size, e = l.Len()
	assert.Nil(t, e)
	assert.Equal(t, size, count+1)
	time.Sleep(duration/2 + duration/3)

	size, e = l.Len()
	assert.Nil(t, e)
	assert.Equal(t, size, 1)

	key := "ok"
	val, e := l.Get(key)
	assert.Nil(t, e)
	assert.Equal(t, val, "value")

	time.Sleep(duration / 2)
	size, e = l.Len()
	assert.Nil(t, e)
	assert.Equal(t, size, 1)
	key = "ok"
	val, e = l.Get(key)
	assert.Nil(t, e)
	assert.Equal(t, val, "value")

	time.Sleep(duration + duration/3)
	size, e = l.Len()
	assert.Nil(t, e)
	assert.Equal(t, size, 0)

	// batch
	e = l.BatchPut(1, "1", 2)
	assert.Nil(t, e)
	size, e = l.Len()
	assert.Nil(t, e)
	assert.Equal(t, size, 2)
	vals, e := l.BatchGet(1, 2, 3)
	assert.Nil(t, e)
	assert.Equal(t, len(vals), 3)
	assert.True(t, vals[0].Exists)
	assert.Equal(t, vals[0].Value, "1")
	assert.True(t, vals[1].Exists)
	assert.Nil(t, vals[1].Value)
	assert.False(t, vals[2].Exists)
	assert.Nil(t, vals[2].Value)
}
