package model

import (
	"strings"
	"time"
)

const (
	SignRequestPending = "pending"
	SignRequestSigned  = "signed"
	SignRequestFailed  = "failed"
)

var signRequestTransitions = map[string]map[string]bool{
	SignRequestPending: {SignRequestSigned: true, SignRequestFailed: true},
}

func SignRequestCanTransition(from, to string) bool {
	if m, ok := signRequestTransitions[from]; ok {
		return m[to]
	}
	return false
}

type SignRequest struct {
	ID          string     `json:"id"`
	KeyPairID   string     `json:"key_pair_id"`
	Algorithm   string     `json:"algorithm"`
	PayloadHash string     `json:"payload_hash"`
	RequestID   string     `json:"request_id"`
	Status      string     `json:"status"`
	SignedAt    *time.Time `json:"signed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (s *SignRequest) Validate() error {
	s.KeyPairID = strings.TrimSpace(s.KeyPairID)
	if s.KeyPairID == "" {
		return NewValidationError("key_pair_id", "密钥对 ID 不能为空")
	}
	s.Algorithm = strings.TrimSpace(s.Algorithm)
	if s.Algorithm == "" {
		return NewValidationError("algorithm", "算法不能为空")
	}
	s.PayloadHash = strings.TrimSpace(s.PayloadHash)
	if s.PayloadHash == "" {
		return NewValidationError("payload_hash", "载荷哈希不能为空")
	}
	s.RequestID = strings.TrimSpace(s.RequestID)
	if s.RequestID == "" {
		return NewValidationError("request_id", "请求 ID 不能为空")
	}
	if s.Status == "" {
		s.Status = SignRequestPending
	}
	if s.Status != SignRequestPending && s.Status != SignRequestSigned && s.Status != SignRequestFailed {
		return NewValidationError("status", "签名请求状态不合法")
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now()
	}
	return nil
}

type SignRequestFilter struct {
	KeyPairID string
	Algorithm string
	Status    string
	RequestID string
}

func (f SignRequestFilter) Match(s *SignRequest) bool {
	if f.KeyPairID != "" && s.KeyPairID != f.KeyPairID {
		return false
	}
	if f.Algorithm != "" && s.Algorithm != f.Algorithm {
		return false
	}
	if f.Status != "" && s.Status != f.Status {
		return false
	}
	if f.RequestID != "" && s.RequestID != f.RequestID {
		return false
	}
	return true
}
