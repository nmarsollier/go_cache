package memoize

import (
	"sync"
	"sync/atomic"
)

type SafeMemoize struct {
	cache   *Memo
	mutex   *sync.Mutex
	loading int32
}

// InvalidateCache invalidates the cache
func (m *SafeMemoize) InvalidateCache() {
	m.cache = nil
}

// ReplaceMockCache just to mock test, this should be removed
func (m *SafeMemoize) ReplaceMockCache(newCache *Memo) {
	m.cache = newCache
}

// Value get cached value, fetching data if needed
func (m *SafeMemoize) Value(
	fetchFunc func() *Memo,
) interface{} {
	currCache := m.cache
	if currCache != nil && currCache.Value() != nil {
		return currCache.Cached()
	}

	loadData := atomic.CompareAndSwapInt32(&m.loading, 0, 1)
	switch {
	case currCache == nil:
		// No data cached, just lock the concurrent calls
		currCache = m.fetchData(fetchFunc)
	case loadData:
		// There is a valid cache, and we need to load data, do it in goroutine
		go m.fetchData(fetchFunc)
	case !loadData:
		// No need to load data, just return the actual value
		return currCache.Cached()
	}

	currCache = m.cache
	return currCache.Cached()
}

func (m *SafeMemoize) fetchData(
	fetchFunc func() *Memo,
) *Memo {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	currCache := m.cache
	if m.loading == 0 && currCache != nil && currCache.Value() != nil {
		return currCache
	}

	defer func() { m.loading = 0 }()

	currCache = fetchFunc()
	m.cache = currCache
	return currCache
}

// NewSafeMemoize creates new thread safe memoization
func NewSafeMemoize() *SafeMemoize {
	return &SafeMemoize{
		cache:   nil,
		mutex:   &sync.Mutex{},
		loading: 0,
	}
}
