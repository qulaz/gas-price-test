package cache

import "time"

//go:generate go run github.com/golang/mock/mockgen -source=interface.go -destination=cachemock/mock.go -package=cachemock

type ExpiredCache[K comparable, V any] interface {
	Get(key K) (V, error)
	Upsert(key K, value V, ttl time.Duration) error
	Delete(key K) error
}
