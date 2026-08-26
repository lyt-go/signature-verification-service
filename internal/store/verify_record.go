package store

import (
	"signatureservice/internal/model"
)

func (s *MemoryStore) CreateVerifyRecord(v *model.VerifyRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.verifyRecords[v.ID] = v.Clone()
	return nil
}

func (s *MemoryStore) GetVerifyRecord(id string) (*model.VerifyRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.verifyRecords[id]
	if !ok {
		return nil, ErrNotFound
	}
	return v.Clone(), nil
}

func (s *MemoryStore) ListVerifyRecords() []*model.VerifyRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.VerifyRecord, 0, len(s.verifyRecords))
	for _, v := range s.verifyRecords {
		list = append(list, v.Clone())
	}
	return list
}

func (s *MemoryStore) UpdateVerifyRecord(v *model.VerifyRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.verifyRecords[v.ID]; !ok {
		return ErrNotFound
	}
	s.verifyRecords[v.ID] = v.Clone()
	return nil
}

func (s *MemoryStore) DeleteVerifyRecord(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.verifyRecords[id]; !ok {
		return ErrNotFound
	}
	delete(s.verifyRecords, id)
	return nil
}

func (s *MemoryStore) BatchCreateVerifyRecords(records []*model.VerifyRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, v := range records {
		s.verifyRecords[v.ID] = v.Clone()
	}
	return nil
}
