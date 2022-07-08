package cache

type memoryCache struct {
	m map[string][]byte
}

func NewMemoryCache() Cache {
	return memoryCache{m: make(map[string][]byte)}
}

func (m memoryCache) Get(key string) ([]byte, bool) {
	data, ok := m.m[key]

	return data, ok
}

func (m memoryCache) Set(key string, data []byte) {
	m.m[key] = data
}
