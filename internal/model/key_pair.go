package model

import (
	"strings"
	"time"
)

const (
	KeyPairActive  = "active"
	KeyPairRevoked = "revoked"
	KeyPairExpired = "expired"
)

var keyPairTransitions = map[string]map[string]bool{
	KeyPairActive: {KeyPairRevoked: true, KeyPairExpired: true},
}

func KeyPairCanTransition(from, to string) bool {
	if m, ok := keyPairTransitions[from]; ok {
		return m[to]
	}
	return false
}

type KeyPair struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Algorithm string    `json:"algorithm"`
	KeySize   int       `json:"key_size"`
	PublicKey string    `json:"public_key"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (k *KeyPair) Validate() error {
	k.Name = strings.TrimSpace(k.Name)
	if k.Name == "" {
		return NewValidationError("name", "密钥对名称不能为空")
	}
	k.Algorithm = strings.TrimSpace(k.Algorithm)
	if k.Algorithm == "" {
		return NewValidationError("algorithm", "算法不能为空")
	}
	if k.Algorithm != "rsa" && k.Algorithm != "ecdsa" && k.Algorithm != "ed25519" {
		return NewValidationError("algorithm", "算法必须是 rsa、ecdsa 或 ed25519")
	}
	if k.KeySize <= 0 {
		return NewValidationError("key_size", "密钥长度必须大于 0")
	}
	if k.Status == "" {
		k.Status = KeyPairActive
	}
	if k.Status != KeyPairActive && k.Status != KeyPairRevoked && k.Status != KeyPairExpired {
		return NewValidationError("status", "密钥对状态不合法")
	}
	if k.CreatedAt.IsZero() {
		k.CreatedAt = time.Now()
	}
	return nil
}

type KeyPairFilter struct {
	Name      string
	Algorithm string
	Status    string
	Keyword   string
}

func (f KeyPairFilter) Match(k *KeyPair) bool {
	if f.Name != "" && k.Name != f.Name {
		return false
	}
	if f.Algorithm != "" && k.Algorithm != f.Algorithm {
		return false
	}
	if f.Status != "" && k.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		kw := strings.ToLower(strings.TrimSpace(f.Keyword))
		if kw != "" && !strings.Contains(strings.ToLower(k.Name), kw) &&
			!strings.Contains(strings.ToLower(k.PublicKey), kw) {
			return false
		}
	}
	return true
}
