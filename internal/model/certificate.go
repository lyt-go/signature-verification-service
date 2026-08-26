package model

import (
	"strings"
	"time"
)

const (
	CertificateValid   = "valid"
	CertificateRevoked = "revoked"
	CertificateExpired = "expired"
)

var certificateTransitions = map[string]map[string]bool{
	CertificateValid: {CertificateRevoked: true, CertificateExpired: true},
}

func CertificateCanTransition(from, to string) bool {
	if m, ok := certificateTransitions[from]; ok {
		return m[to]
	}
	return false
}

type Certificate struct {
	ID           string    `json:"id"`
	KeyPairID    string    `json:"key_pair_id"`
	Subject      string    `json:"subject"`
	Issuer       string    `json:"issuer"`
	SerialNumber string    `json:"serial_number"`
	Status       string    `json:"status"`
	ValidFrom    time.Time `json:"valid_from"`
	ValidTo      time.Time `json:"valid_to"`
}

func (c *Certificate) Validate() error {
	c.KeyPairID = strings.TrimSpace(c.KeyPairID)
	if c.KeyPairID == "" {
		return NewValidationError("key_pair_id", "密钥对 ID 不能为空")
	}
	c.Subject = strings.TrimSpace(c.Subject)
	if c.Subject == "" {
		return NewValidationError("subject", "主题不能为空")
	}
	c.Issuer = strings.TrimSpace(c.Issuer)
	if c.Issuer == "" {
		return NewValidationError("issuer", "颁发者不能为空")
	}
	c.SerialNumber = strings.TrimSpace(c.SerialNumber)
	if c.SerialNumber == "" {
		return NewValidationError("serial_number", "序列号不能为空")
	}
	if c.Status == "" {
		c.Status = CertificateValid
	}
	if c.Status != CertificateValid && c.Status != CertificateRevoked && c.Status != CertificateExpired {
		return NewValidationError("status", "证书状态不合法")
	}
	if c.ValidFrom.IsZero() {
		c.ValidFrom = time.Now()
	}
	if c.ValidTo.IsZero() {
		c.ValidTo = c.ValidFrom.Add(365 * 24 * time.Hour)
	}
	if !c.ValidTo.After(c.ValidFrom) {
		return NewValidationError("valid_to", "有效期截止时间必须晚于开始时间")
	}
	return nil
}

type CertificateFilter struct {
	KeyPairID string
	Status    string
	Subject   string
	Issuer    string
}

func (f CertificateFilter) Match(c *Certificate) bool {
	if f.KeyPairID != "" && c.KeyPairID != f.KeyPairID {
		return false
	}
	if f.Status != "" && c.Status != f.Status {
		return false
	}
	if f.Subject != "" && c.Subject != f.Subject {
		return false
	}
	if f.Issuer != "" && c.Issuer != f.Issuer {
		return false
	}
	return true
}
