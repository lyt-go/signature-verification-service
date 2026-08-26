// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"signatureservice/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// KeyPair
	CreateKeyPair(k *model.KeyPair) error
	GetKeyPair(id string) (*model.KeyPair, error)
	ListKeyPairs() []*model.KeyPair
	UpdateKeyPair(k *model.KeyPair) error
	DeleteKeyPair(id string) error

	// SignRequest
	CreateSignRequest(sr *model.SignRequest) error
	GetSignRequest(id string) (*model.SignRequest, error)
	ListSignRequests() []*model.SignRequest
	UpdateSignRequest(sr *model.SignRequest) error
	DeleteSignRequest(id string) error

	// Signature
	CreateSignature(sig *model.Signature) error
	GetSignature(id string) (*model.Signature, error)
	GetSignatureBySignRequestID(signRequestID string) (*model.Signature, error)
	ListSignatures() []*model.Signature
	UpdateSignature(sig *model.Signature) error
	DeleteSignature(id string) error

	// VerifyRecord
	CreateVerifyRecord(v *model.VerifyRecord) error
	GetVerifyRecord(id string) (*model.VerifyRecord, error)
	ListVerifyRecords() []*model.VerifyRecord
	UpdateVerifyRecord(v *model.VerifyRecord) error
	DeleteVerifyRecord(id string) error
	BatchCreateVerifyRecords(records []*model.VerifyRecord) error

	// Algorithm
	CreateAlgorithm(a *model.Algorithm) error
	GetAlgorithm(id string) (*model.Algorithm, error)
	GetAlgorithmByName(name string) (*model.Algorithm, error)
	ListAlgorithms() []*model.Algorithm
	UpdateAlgorithm(a *model.Algorithm) error
	DeleteAlgorithm(id string) error

	// Certificate
	CreateCertificate(c *model.Certificate) error
	GetCertificate(id string) (*model.Certificate, error)
	ListCertificates() []*model.Certificate
	UpdateCertificate(c *model.Certificate) error
	DeleteCertificate(id string) error
}
