package store

import (
	"signatureservice/internal/model"
	"sync"
)

type VerificationJobStore struct {
	mu   sync.Mutex
	jobs map[string]*model.VerificationJob
}

func NewVerificationJobStore() *VerificationJobStore {
	return &VerificationJobStore{jobs: make(map[string]*model.VerificationJob)}
}
func (s *VerificationJobStore) Save(job *model.VerificationJob) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	copy := *job
	s.jobs[job.ID] = &copy
	return nil
}
func (s *VerificationJobStore) Get(id string) *model.VerificationJob {
	s.mu.Lock()
	defer s.mu.Unlock()
	job := *s.jobs[id]
	return &job
}
