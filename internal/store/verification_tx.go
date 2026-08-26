package store

import (
	"errors"
	"sync"
)

var ErrTemporaryCommit = errors.New("temporary commit failure")

type VerificationTxStore struct {
	mu          sync.Mutex
	failCommits int
	records     []string
}

func NewVerificationTxStore(failCommits int) *VerificationTxStore {
	return &VerificationTxStore{failCommits: failCommits}
}

type VerificationTx struct {
	store *VerificationTxStore
	id    string
}

func (s *VerificationTxStore) Begin(id string) *VerificationTx {
	return &VerificationTx{store: s, id: id}
}

func (tx *VerificationTx) Commit() error {
	tx.store.mu.Lock()
	defer tx.store.mu.Unlock()
	tx.store.records = append(tx.store.records, tx.id)
	if tx.store.failCommits > 0 {
		tx.store.failCommits--
		return ErrTemporaryCommit
	}
	return nil
}

func (tx *VerificationTx) Rollback() error { return nil }
func (s *VerificationTxStore) RecordCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.records)
}
