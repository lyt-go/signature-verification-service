package model

import (
	"strings"
	"time"
)

type Signature struct {
	ID            string    `json:"id"`
	SignRequestID string    `json:"sign_request_id"`
	KeyPairID     string    `json:"key_pair_id"`
	Value         string    `json:"value"`
	Algorithm     string    `json:"algorithm"`
	CreatedAt     time.Time `json:"created_at"`
}

func (s *Signature) Validate() error {
	s.SignRequestID = strings.TrimSpace(s.SignRequestID)
	if s.SignRequestID == "" {
		return NewValidationError("sign_request_id", "签名请求 ID 不能为空")
	}
	s.KeyPairID = strings.TrimSpace(s.KeyPairID)
	if s.KeyPairID == "" {
		return NewValidationError("key_pair_id", "密钥对 ID 不能为空")
	}
	s.Value = strings.TrimSpace(s.Value)
	if s.Value == "" {
		return NewValidationError("value", "签名值不能为空")
	}
	s.Algorithm = strings.TrimSpace(s.Algorithm)
	if s.Algorithm == "" {
		return NewValidationError("algorithm", "算法不能为空")
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now()
	}
	return nil
}

// Clone 返回 s 的深拷贝，用于 store 读写边界上的快照隔离。
func (s *Signature) Clone() *Signature {
	cp := *s
	return &cp
}

type SignatureFilter struct {
	KeyPairID     string
	Algorithm     string
	SignRequestID string
}

func (f SignatureFilter) Match(s *Signature) bool {
	if f.KeyPairID != "" && s.KeyPairID != f.KeyPairID {
		return false
	}
	if f.Algorithm != "" && s.Algorithm != f.Algorithm {
		return false
	}
	if f.SignRequestID != "" && s.SignRequestID != f.SignRequestID {
		return false
	}
	return true
}
