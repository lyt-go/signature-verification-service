package store

import (
	"sync"

	"signatureservice/internal/model"
)

type MemoryStore struct {
	mu            sync.RWMutex
	keyPairs      map[string]*model.KeyPair
	signRequests  map[string]*model.SignRequest
	signatures    map[string]*model.Signature
	verifyRecords map[string]*model.VerifyRecord
	algorithms    map[string]*model.Algorithm
	certificates  map[string]*model.Certificate
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		keyPairs:      make(map[string]*model.KeyPair),
		signRequests:  make(map[string]*model.SignRequest),
		signatures:    make(map[string]*model.Signature),
		verifyRecords: make(map[string]*model.VerifyRecord),
		algorithms:    make(map[string]*model.Algorithm),
		certificates:  make(map[string]*model.Certificate),
	}
}

var _ Store = (*MemoryStore)(nil)
