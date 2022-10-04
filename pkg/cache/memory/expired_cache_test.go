package memory

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/qulaz/gas-price-test/pkg/cache"
)

func TestExpiredCache_Upsert(t *testing.T) {
	c := NewExpiredCache[string, string]()

	t.Run("insert first value", func(t *testing.T) {
		key, value := "first", "value"
		err := c.Upsert(key, value, time.Millisecond)
		require.NoError(t, err)
		assert.Equal(t, c.cache[key].value, value)
		assert.Len(t, c.cache, 1)
	})
	t.Run("replace first value", func(t *testing.T) {
		key, value := "first", "value1"
		err := c.Upsert(key, value, time.Millisecond)
		require.NoError(t, err)
		assert.Equal(t, c.cache[key].value, value)
		assert.Len(t, c.cache, 1)
	})
	t.Run("insert second value", func(t *testing.T) {
		key, value := "second", "value"
		err := c.Upsert(key, value, time.Millisecond)
		require.NoError(t, err)
		assert.Equal(t, c.cache[key].value, value)
		assert.Len(t, c.cache, 2)
	})
	t.Run("concurrency", func(t *testing.T) {
		c := NewExpiredCache[int, int]()

		var wg sync.WaitGroup
		kv := map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6, 7: 7, 8: 8, 9: 9, 10: 10, 11: 11}

		for key, value := range kv {
			wg.Add(1)
			go func(key, value int) {
				defer wg.Done()

				err := c.Upsert(key, value, time.Millisecond)
				require.NoError(t, err)
			}(key, value)
		}

		wg.Wait()

		require.Len(t, c.cache, len(kv))

		for k, v := range c.cache {
			assert.Equal(t, kv[k], v.value)
		}
	})
}

func TestExpiredCache_Get(t *testing.T) {
	c := NewExpiredCache[int, int]()

	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	start := time.Now().UTC()
	ttl := time.Millisecond

	for _, v := range values {
		err := c.Upsert(v, v, time.Millisecond)
		require.NoError(t, err)
	}

	t.Run("existed key", func(t *testing.T) {
		if time.Since(start) > ttl-ttl/2 {
			t.Log("key already expired, try increase ttl")
			t.SkipNow()
		}

		expectedValue := values[0]
		v, err := c.Get(expectedValue)
		require.NoError(t, err)
		assert.Equal(t, expectedValue, v)
	})
	t.Run("not found key", func(t *testing.T) {
		_, err := c.Get(3425)
		require.Error(t, err)
		require.True(t, errors.Is(err, cache.ErrKeyNotFound))
		assert.Len(t, c.cache, len(values))
	})
	t.Run("expired key", func(t *testing.T) {
		time.Sleep(ttl)

		_, err := c.Get(values[0])
		require.Error(t, err)
		require.True(t, errors.Is(err, cache.ErrKeyNotFound))
		assert.Len(t, c.cache, len(values)-1)
	})
	t.Run("concurrency", func(t *testing.T) {
		c := NewExpiredCache[int, int]()
		values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

		for _, v := range values {
			err := c.Upsert(v, v, time.Second*5)
			require.NoError(t, err)
		}

		var wg sync.WaitGroup

		for _, key := range values {
			wg.Add(1)
			go func(key int) {
				defer wg.Done()

				value, err := c.Get(key)
				require.NoError(t, err)
				assert.Equal(t, key, value)
			}(key)
		}

		wg.Wait()
	})
}

func TestExpiredCache_Delete(t *testing.T) {
	c := NewExpiredCache[int, int]()

	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	for _, v := range values {
		err := c.Upsert(v, v, time.Second)
		require.NoError(t, err)
	}

	t.Run("success", func(t *testing.T) {
		err := c.Delete(values[0])
		require.NoError(t, err)
		assert.Len(t, c.cache, len(values)-1)
	})
	t.Run("key not found", func(t *testing.T) {
		err := c.Delete(35235)
		require.Error(t, err)
		assert.True(t, errors.Is(err, cache.ErrKeyNotFound))
	})
	t.Run("concurrency", func(t *testing.T) {
		c := NewExpiredCache[int, int]()
		values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

		for _, v := range values {
			err := c.Upsert(v, v, time.Second*5)
			require.NoError(t, err)
		}

		var wg sync.WaitGroup

		for _, key := range values {
			wg.Add(1)
			go func(key int) {
				defer wg.Done()

				err := c.Delete(key)
				require.NoError(t, err)
			}(key)
		}

		wg.Wait()

		assert.Len(t, c.cache, 0)
	})
}
