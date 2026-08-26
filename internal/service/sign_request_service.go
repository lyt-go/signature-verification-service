package service

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"time"

	"signatureservice/internal/model"
	"signatureservice/pkg/idgen"
)

func (s *Service) CreateSignRequest(input model.SignRequest) (*model.SignRequest, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	input.ID = idgen.Hex()
	if err := s.store.CreateSignRequest(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) ListSignRequests(filter model.SignRequestFilter, page, size int) ([]*model.SignRequest, int, error) {
	all := s.store.ListSignRequests()
	matched := make([]*model.SignRequest, 0, len(all))
	for _, sr := range all {
		if filter.Match(sr) {
			matched = append(matched, sr)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.SignRequest{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetSignRequest(id string) (*model.SignRequest, error) {
	return s.store.GetSignRequest(id)
}

func (s *Service) UpdateSignRequest(id string, input model.SignRequest) (*model.SignRequest, error) {
	sr, err := s.store.GetSignRequest(id)
	if err != nil {
		return nil, err
	}
	if input.Algorithm != "" {
		sr.Algorithm = input.Algorithm
	}
	if input.PayloadHash != "" {
		sr.PayloadHash = input.PayloadHash
	}
	if input.RequestID != "" {
		sr.RequestID = input.RequestID
	}
	if err := sr.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateSignRequest(sr); err != nil {
		return nil, err
	}
	return sr, nil
}

func (s *Service) UpdateSignRequestStatus(id string, newStatus string) (*model.SignRequest, error) {
	sr, err := s.store.GetSignRequest(id)
	if err != nil {
		return nil, err
	}
	if !model.SignRequestCanTransition(sr.Status, newStatus) {
		return nil, model.NewValidationError("status", "状态流转不合法")
	}
	sr.Status = newStatus
	now := time.Now()
	if newStatus == model.SignRequestSigned {
		sr.SignedAt = &now
	}
	if err := s.store.UpdateSignRequest(sr); err != nil {
		return nil, err
	}
	return sr, nil
}

func (s *Service) DeleteSignRequest(id string) error {
	return s.store.DeleteSignRequest(id)
}

func (s *Service) ProcessSignRequest(id string) (*model.Signature, error) {
	sr, err := s.store.GetSignRequest(id)
	if err != nil {
		return nil, err
	}
	if sr.Status != model.SignRequestPending {
		return nil, model.NewValidationError("status", "只有 pending 状态的请求才能执行签名")
	}
	kp, err := s.store.GetKeyPair(sr.KeyPairID)
	if err != nil {
		return nil, model.NewValidationError("key_pair_id", "关联的密钥对不存在")
	}
	if kp.Status != model.KeyPairActive {
		return nil, model.NewValidationError("key_pair_id", "关联的密钥对状态不是 active，无法签名")
	}
	h := sha256.Sum256([]byte(sr.PayloadHash))
	sigValue := hex.EncodeToString(h[:])

	now := time.Now()
	sr.Status = model.SignRequestSigned
	sr.SignedAt = &now
	if err := s.store.UpdateSignRequest(sr); err != nil {
		return nil, err
	}

	sig := &model.Signature{
		ID:            idgen.Hex(),
		SignRequestID: sr.ID,
		KeyPairID:     sr.KeyPairID,
		Value:         sigValue,
		Algorithm:     sr.Algorithm,
		CreatedAt:     now,
	}
	if err := sig.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateSignature(sig); err != nil {
		return nil, err
	}
	return sig, nil
}
