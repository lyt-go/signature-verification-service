package store

import "sync"

type VerificationBatchCache struct {
	mu       sync.RWMutex
	payloads map[string][]byte
}

func NewVerificationBatchCache() *VerificationBatchCache {
	return &VerificationBatchCache{payloads: make(map[string][]byte)}
}
func (c *VerificationBatchCache) Save(id string, payload []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.payloads[id] = payload
}
func (c *VerificationBatchCache) Get(id string) []byte {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.payloads[id]
}
