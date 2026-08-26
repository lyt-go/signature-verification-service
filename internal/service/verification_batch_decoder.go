package service

import "signatureservice/internal/store"

type VerificationBatchDecoder struct {
	cache *store.VerificationBatchCache
}

func NewVerificationBatchDecoder(cache *store.VerificationBatchCache) *VerificationBatchDecoder {
	return &VerificationBatchDecoder{cache: cache}
}
func (d *VerificationBatchDecoder) DecodeAndExport(id string, raw []byte, release <-chan struct{}, exported chan<- string) []byte {
	payload := make([]byte, len(raw))
	copy(payload, raw)
	d.cache.Save(id, payload)
	go func(value []byte) { <-release; exported <- string(value) }(payload)
	return payload
}
