package store

import (
	"signatureservice/internal/model"
)

func (s *MemoryStore) CreateSignRequest(sr *model.SignRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.signRequests[sr.ID] = sr
	return nil
}

func (s *MemoryStore) GetSignRequest(id string) (*model.SignRequest, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sr, ok := s.signRequests[id]
	if !ok {
		return nil, ErrNotFound
	}
	return sr, nil
}

func (s *MemoryStore) ListSignRequests() []*model.SignRequest {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.SignRequest, 0, len(s.signRequests))
	for _, sr := range s.signRequests {
		list = append(list, sr)
	}
	return list
}

func (s *MemoryStore) UpdateSignRequest(sr *model.SignRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.signRequests[sr.ID]; !ok {
		return ErrNotFound
	}
	s.signRequests[sr.ID] = sr
	return nil
}

func (s *MemoryStore) DeleteSignRequest(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.signRequests[id]; !ok {
		return ErrNotFound
	}
	delete(s.signRequests, id)
	return nil
}
