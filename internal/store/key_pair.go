package store

import (
	"signatureservice/internal/model"
)

func (s *MemoryStore) CreateKeyPair(k *model.KeyPair) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.keyPairs[k.ID] = k
	return nil
}

func (s *MemoryStore) GetKeyPair(id string) (*model.KeyPair, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	k, ok := s.keyPairs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return k, nil
}

func (s *MemoryStore) ListKeyPairs() []*model.KeyPair {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.KeyPair, 0, len(s.keyPairs))
	for _, k := range s.keyPairs {
		list = append(list, k)
	}
	return list
}

func (s *MemoryStore) UpdateKeyPair(k *model.KeyPair) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.keyPairs[k.ID]; !ok {
		return ErrNotFound
	}
	s.keyPairs[k.ID] = k
	return nil
}

func (s *MemoryStore) DeleteKeyPair(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.keyPairs[id]; !ok {
		return ErrNotFound
	}
	delete(s.keyPairs, id)
	return nil
}
