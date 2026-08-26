package store

import (
	"fmt"
	"signatureservice/internal/model"
	"sync"
)

type CertificateBundle struct {
	ID          string
	Certificate *model.Certificate
	Signature   *model.Signature
}
type BundleCache struct {
	mu      sync.RWMutex
	bundles map[string]*CertificateBundle
}

func NewBundleCache() *BundleCache { return &BundleCache{bundles: make(map[string]*CertificateBundle)} }
func (c *BundleCache) Put(bundle *CertificateBundle) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.bundles[bundle.ID] = bundle
}
func (c *BundleCache) Get(id string) (*CertificateBundle, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	bundle, ok := c.bundles[id]
	if !ok {
		return nil, ErrNotFound
	}
	// 拒绝读取残缺包：未校验通过的组合包即便因历史 bug 被写入缓存，
	// 这里也不会返回，更不会因为字段为 nil 而触发 panic。
	if err := ValidateBundle(bundle); err != nil {
		return nil, err
	}
	return bundle, nil
}
func ValidateBundle(bundle *CertificateBundle) error {
	if bundle.Certificate == nil || bundle.Signature == nil {
		return fmt.Errorf("incomplete certificate bundle")
	}
	return nil
}
