package cache

import (
	"context"
	"sync"
	"time"
)

type cacheEntity[T any] struct {
	data T
	time time.Time
}

func newCacheEntity[T any](data T, ttl time.Duration) cacheEntity[T] {
	return cacheEntity[T]{
		data: data,
		time: time.Now().Add(ttl),
	}
}

type Cache[K comparable, T any] struct {
	ctx context.Context
	ttl time.Duration

	mux  sync.Mutex
	data map[K]cacheEntity[T]
}

func New[K comparable, T any](ctx context.Context, ttl time.Duration) *Cache[K, T] {
	cache := &Cache[K, T]{
		data: make(map[K]cacheEntity[T]),
		ctx:  ctx,
		ttl:  ttl,
	}

	go cache.cleanCache()

	return cache
}

func (c *Cache[K, T]) Get(k K) (T, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	cacheEntity, ok := c.data[k]
	if !ok {
		return cacheEntity.data, false
	}

	cacheEntity.time = time.Now().Add(c.ttl)
	c.data[k] = cacheEntity

	return cacheEntity.data, true
}

func (c *Cache[K, T]) Set(k K, v T) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.data[k] = newCacheEntity(v, c.ttl)
}

// cleanCache инвалидирует данные в кеше, если они устарели
// Работает по схеме crawler
func (c *Cache[K, T]) cleanCache() {
	ticker := time.NewTicker(c.ttl / 2)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.mux.Lock()
			for key, entity := range c.data {
				if time.Now().After(entity.time) {
					delete(c.data, key)
				}
			}
			c.mux.Unlock()
		}
	}
}
