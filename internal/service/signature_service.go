package service

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"

	"signatureservice/internal/model"
	"signatureservice/pkg/idgen"
)

func (s *Service) CreateSignature(input model.Signature) (*model.Signature, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	input.ID = idgen.Hex()
	if err := s.store.CreateSignature(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) ListSignatures(filter model.SignatureFilter, page, size int) ([]*model.Signature, int, error) {
	all := s.store.ListSignatures()
	matched := make([]*model.Signature, 0, len(all))
	for _, sig := range all {
		if filter.Match(sig) {
			matched = append(matched, sig)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Signature{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetSignature(id string) (*model.Signature, error) {
	return s.store.GetSignature(id)
}

func (s *Service) UpdateSignature(id string, input model.Signature) (*model.Signature, error) {
	sig, err := s.store.GetSignature(id)
	if err != nil {
		return nil, err
	}
	if input.Value != "" {
		sig.Value = input.Value
	}
	if input.Algorithm != "" {
		sig.Algorithm = input.Algorithm
	}
	if err := sig.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateSignature(sig); err != nil {
		return nil, err
	}
	return sig, nil
}

func (s *Service) DeleteSignature(id string) error {
	return s.store.DeleteSignature(id)
}

func (s *Service) VerifySignature(signatureID, payloadHash, verifier string) (*model.VerifyRecord, error) {
	sig, err := s.store.GetSignature(signatureID)
	if err != nil {
		return nil, err
	}
	h := sha256.Sum256([]byte(payloadHash))
	computed := hex.EncodeToString(h[:])
	valid := computed == sig.Value

	vr := &model.VerifyRecord{
		ID:          idgen.Hex(),
		SignatureID: signatureID,
		PayloadHash: payloadHash,
		Valid:       valid,
		Verifier:    verifier,
	}
	if err := vr.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateVerifyRecord(vr); err != nil {
		return nil, err
	}
	return vr, nil
}
