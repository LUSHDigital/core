package keys

import (
	"context"
	"sync"
)

// NewCache returns a cache for sources
func NewCache(source Source) Source {
	return &Cache{
		source: source,
		lock:   &sync.Mutex{},
	}
}

// Cache represents a cache of sources
type Cache struct {
	cached    bool
	cachedKey []byte
	cachedErr error
	source    Source
	lock      sync.Locker
}

// GetKey resolves the cache source and caches it
func (c Cache) GetKey(ctx context.Context) ([]byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if !c.cached {
		c.cachedKey, c.cachedErr = c.source.GetKey(ctx)
	}
	return c.cachedKey, c.cachedErr
}
