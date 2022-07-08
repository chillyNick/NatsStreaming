package cache

import "sync"

type memoryCache struct {
	m       map[string][]byte
	rwMutex *sync.RWMutex
}

func NewMemoryCache() Cache {
	return memoryCache{m: make(map[string][]byte), rwMutex: &sync.RWMutex{}}
}

func (m memoryCache) Get(key string) ([]byte, bool) {
	m.rwMutex.RLock()
	data, ok := m.m[key]
	m.rwMutex.RUnlock()

	return data, ok
}

func (m memoryCache) Set(key string, data []byte) {
	m.rwMutex.Lock()
	m.m[key] = data
	m.rwMutex.Unlock()
}
