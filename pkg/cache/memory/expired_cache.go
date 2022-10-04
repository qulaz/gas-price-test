package memory

import (
	"sync"
	"time"

	"github.com/qulaz/gas-price-test/pkg/cache"
)

var _ cache.ExpiredCache[string, string] = (*ExpiredCache[string, string])(nil)

type expireValue[V any] struct {
	value    V
	expireAt time.Time
}

// ExpiredCache кеш в памяти на основе map и RWMutex, в котором ключи живут заданное количество
// времени. Из особенностей реализации стоит отметить, что ключи удаляются только при их запросе.
// Если после истечения срока жизни ключа за ним никто не обратится, он так и останется в мапе
// и будет занимать память.
type ExpiredCache[K comparable, V any] struct {
	mutex sync.RWMutex
	cache map[K]expireValue[V]
}

func NewExpiredCache[K comparable, V any]() *ExpiredCache[K, V] {
	return &ExpiredCache[K, V]{
		mutex: sync.RWMutex{},
		cache: make(map[K]expireValue[V]),
	}
}

func (m *ExpiredCache[K, V]) Get(key K) (V, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var noop V

	value, ok := m.cache[key]
	if !ok {
		return noop, cache.ErrKeyNotFound
	}

	if time.Now().UTC().After(value.expireAt) {
		delete(m.cache, key)

		return noop, cache.ErrKeyNotFound
	}

	return value.value, nil
}

func (m *ExpiredCache[K, V]) Upsert(key K, value V, ttl time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.cache[key] = expireValue[V]{
		value:    value,
		expireAt: time.Now().UTC().Add(ttl),
	}

	return nil
}

func (m *ExpiredCache[K, V]) Delete(key K) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, ok := m.cache[key]
	if !ok {
		return cache.ErrKeyNotFound
	}

	delete(m.cache, key)

	return nil
}
