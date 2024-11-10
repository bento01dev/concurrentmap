package concurrentmap

import "sync"

var SHARD_COUNT = 32

type shard[K comparable, V any] struct {
	items map[K]V
	mu    sync.RWMutex
}

type HashFunc[K comparable] func(key K) uint

type ConcurrentMap[K comparable, V any] struct {
	shards   []*shard[K, V]
	hashFunc HashFunc[K]
}

func NewConcurrentMap[K comparable, V any](hashFunc HashFunc[K]) *ConcurrentMap[K, V] {
	m := &ConcurrentMap[K, V]{
		shards:   make([]*shard[K, V], SHARD_COUNT),
		hashFunc: hashFunc,
	}

	for i := 0; i < SHARD_COUNT; i++ {
		m.shards[i] = &shard[K, V]{items: make(map[K]V)}
	}
	return m
}

func (m *ConcurrentMap[K, V]) shardIndex(key K) uint {
	return uint(m.hashFunc(key)) % uint(SHARD_COUNT)
}

func (m *ConcurrentMap[K, V]) Get(key K) (V, bool) {
	shard := m.shards[m.shardIndex(key)]
	shard.mu.RLock()
	defer shard.mu.RUnlock()
	res, ok := shard.items[key]
	return res, ok
}

func (m *ConcurrentMap[K, V]) Set(key K, value V) {
	shard := m.shards[m.shardIndex(key)]
	shard.mu.Lock()
	defer shard.mu.Unlock()
	shard.items[key] = value
}

func (m *ConcurrentMap[K, V]) Remove(key K) {
	shard := m.shards[m.shardIndex(key)]
	shard.mu.Lock()
	defer shard.mu.Unlock()
	delete(shard.items, key)
}

func (m *ConcurrentMap[K, V]) Count() int {
	var count int
	for _, shard := range m.shards {
		shard.mu.RLock()
		count = count + len(shard.items)
		shard.mu.RUnlock()
	}
	return count
}
