package model

import (
	"strings"
	"time"
)

type VerifyRecord struct {
	ID          string    `json:"id"`
	SignatureID string    `json:"signature_id"`
	PayloadHash string    `json:"payload_hash"`
	Valid       bool      `json:"valid"`
	VerifiedAt  time.Time `json:"verified_at"`
	Verifier    string    `json:"verifier"`
}

func (v *VerifyRecord) Validate() error {
	v.SignatureID = strings.TrimSpace(v.SignatureID)
	if v.SignatureID == "" {
		return NewValidationError("signature_id", "签名 ID 不能为空")
	}
	v.PayloadHash = strings.TrimSpace(v.PayloadHash)
	if v.PayloadHash == "" {
		return NewValidationError("payload_hash", "载荷哈希不能为空")
	}
	v.Verifier = strings.TrimSpace(v.Verifier)
	if v.Verifier == "" {
		return NewValidationError("verifier", "验证者不能为空")
	}
	if v.VerifiedAt.IsZero() {
		v.VerifiedAt = time.Now()
	}
	return nil
}

// Clone 返回 v 的深拷贝，用于 store 读写边界上的快照隔离。
func (v *VerifyRecord) Clone() *VerifyRecord {
	cp := *v
	return &cp
}

type VerifyRecordFilter struct {
	SignatureID string
	Valid       *bool
	Verifier    string
}

func (f VerifyRecordFilter) Match(v *VerifyRecord) bool {
	if f.SignatureID != "" && v.SignatureID != f.SignatureID {
		return false
	}
	if f.Valid != nil && v.Valid != *f.Valid {
		return false
	}
	if f.Verifier != "" && v.Verifier != f.Verifier {
		return false
	}
	return true
}
