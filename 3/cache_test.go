package cache

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestCacheCleanerParallel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ttl := 100 * time.Millisecond
	cache := New[string, int](ctx, ttl)

	cache.Set("key1", 42)
	if val, ok := cache.Get("key1"); ok && val != 42 {
		t.Errorf("Expected 42, got %d", val)
	}

	var wg sync.WaitGroup
	const numWorkers = 100
	const numKeys = 10

	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func(i int) {
			defer wg.Done()
			key := strconv.Itoa(i)
			cache.Set(key, i)
			val, _ := cache.Get(key)
			if val != i {
				t.Errorf("Expected %d for key %s, got %d", i, key, val)
			}
		}(i)
	}

	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func(i int) {
			defer wg.Done()
			key := string(rune('a' + i%numKeys))
			cache.Get(key) // Just test we can read concurrently without panics
		}(i)
	}

	wg.Wait()

	time.Sleep(ttl * 2)
	cache.mux.Lock()
	if len(cache.data) != 0 {
		t.Errorf("Expected cache to be empty after expiration, got %d items", len(cache.data))
	}
	cache.mux.Unlock()

	cache.Set("key2", 100)
	cancel()
	time.Sleep(ttl)
}

func TestCacheCleanerRace(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ttl := 50 * time.Millisecond
	cache := New[string, int](ctx, ttl)

	var wg sync.WaitGroup
	const numOps = 1000

	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < numOps; i++ {
			cache.Set("key", i)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < numOps; i++ {
			cache.Get("key")
		}
	}()

	wg.Wait()
}
