package store

import (
	"errors"

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

// Save 是无条件的全量覆盖写，供主流程顺序写入使用。
func (s *VerificationJobStore) Save(job *model.VerificationJob) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	copy := *job
	s.jobs[job.ID] = &copy
	return nil
}

// SaveIfVersion 仅当任务当前版本等于 expectedVersion 时才落库，返回是否写入成功。
// 迟到的延迟回调（如旧版本的重试回调）会因版本已前进而被拒绝，
// 从而避免把已经 succeeded 的状态覆盖回 running。
func (s *VerificationJobStore) SaveIfVersion(job *model.VerificationJob, expectedVersion int) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	current := s.jobs[job.ID]
	if current != nil && current.Version != expectedVersion {
		return false, nil
	}
	copy := *job
	s.jobs[job.ID] = &copy
	return true, nil
}

func (s *VerificationJobStore) Get(id string) *model.VerificationJob {
	s.mu.Lock()
	defer s.mu.Unlock()
	job := *s.jobs[id]
	return &job
}

// ErrStaleCallback 表示回调基于一个已过期的任务版本，被拒绝执行。
var ErrStaleCallback = errors.New("stale callback: task version already advanced")
