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

// Commit 持久化事务。遇到临时错误时不写入任何记录，确保失败的事务不留痕迹，
// 重试时不会产生重复的已提交记录。
func (tx *VerificationTx) Commit() error {
	tx.store.mu.Lock()
	defer tx.store.mu.Unlock()
	if tx.store.failCommits > 0 {
		tx.store.failCommits--
		return ErrTemporaryCommit
	}
	tx.store.records = append(tx.store.records, tx.id)
	return nil
}

func (tx *VerificationTx) Rollback() error { return nil }
func (s *VerificationTxStore) RecordCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.records)
}
