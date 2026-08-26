package store

import (
	"sync"

	"signatureservice/internal/model"
)

// MemoryStore 是基于 map 的内存存储实现。
//
// 快照语义：所有读方法（Get/List 等）返回的都是内部对象的深拷贝，所有写方法
//（Create/Update 等）也会先拷贝入参再保存。这样调用方拿到的指针始终是当时的
// 独立快照——不会被后续更新篡改，反向也不会把对快照的改动泄漏进存储，从而
// 消除了并发下共享同一指针引发的读写冲突。
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
