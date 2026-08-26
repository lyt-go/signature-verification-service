package store

import (
	"signatureservice/internal/model"
)

func (s *MemoryStore) CreateSignature(sig *model.Signature) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.signatures[sig.ID] = sig.Clone()
	return nil
}

func (s *MemoryStore) GetSignature(id string) (*model.Signature, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sig, ok := s.signatures[id]
	if !ok {
		return nil, ErrNotFound
	}
	return sig.Clone(), nil
}

func (s *MemoryStore) GetSignatureBySignRequestID(signRequestID string) (*model.Signature, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, sig := range s.signatures {
		if sig.SignRequestID == signRequestID {
			return sig.Clone(), nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListSignatures() []*model.Signature {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Signature, 0, len(s.signatures))
	for _, sig := range s.signatures {
		list = append(list, sig.Clone())
	}
	return list
}

func (s *MemoryStore) UpdateSignature(sig *model.Signature) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.signatures[sig.ID]; !ok {
		return ErrNotFound
	}
	s.signatures[sig.ID] = sig.Clone()
	return nil
}

func (s *MemoryStore) DeleteSignature(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.signatures[id]; !ok {
		return ErrNotFound
	}
	delete(s.signatures, id)
	return nil
}
