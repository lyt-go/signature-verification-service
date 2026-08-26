package model

import (
	"strings"
	"time"
)

type Algorithm struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	KeySize   int       `json:"key_size"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}

func (a *Algorithm) Validate() error {
	a.Name = strings.TrimSpace(a.Name)
	if a.Name == "" {
		return NewValidationError("name", "算法名称不能为空")
	}
	a.Type = strings.TrimSpace(a.Type)
	if a.Type == "" {
		return NewValidationError("type", "算法类型不能为空")
	}
	if a.KeySize <= 0 {
		return NewValidationError("key_size", "密钥长度必须大于 0")
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = time.Now()
	}
	return nil
}

type AlgorithmFilter struct {
	Name    string
	Type    string
	Enabled *bool
}

func (f AlgorithmFilter) Match(a *Algorithm) bool {
	if f.Name != "" && a.Name != f.Name {
		return false
	}
	if f.Type != "" && a.Type != f.Type {
		return false
	}
	if f.Enabled != nil && a.Enabled != *f.Enabled {
		return false
	}
	return true
}
