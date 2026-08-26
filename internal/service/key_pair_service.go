package service

import (
	"sort"
	"time"

	"signatureservice/internal/model"
	"signatureservice/pkg/idgen"
)

func (s *Service) CreateKeyPair(input model.KeyPair) (*model.KeyPair, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	input.ID = idgen.Hex()
	if input.ExpiresAt.IsZero() {
		input.ExpiresAt = time.Now().Add(365 * 24 * time.Hour)
	}
	if err := s.store.CreateKeyPair(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) ListKeyPairs(filter model.KeyPairFilter, page, size int) ([]*model.KeyPair, int, error) {
	all := s.store.ListKeyPairs()
	matched := make([]*model.KeyPair, 0, len(all))
	for _, k := range all {
		if filter.Match(k) {
			matched = append(matched, k)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.KeyPair{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetKeyPair(id string) (*model.KeyPair, error) {
	return s.store.GetKeyPair(id)
}

func (s *Service) UpdateKeyPair(id string, input model.KeyPair) (*model.KeyPair, error) {
	k, err := s.store.GetKeyPair(id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		k.Name = input.Name
	}
	if input.PublicKey != "" {
		k.PublicKey = input.PublicKey
	}
	if input.KeySize > 0 {
		k.KeySize = input.KeySize
	}
	if err := k.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateKeyPair(k); err != nil {
		return nil, err
	}
	return k, nil
}

func (s *Service) UpdateKeyPairStatus(id string, newStatus string) (*model.KeyPair, error) {
	k, err := s.store.GetKeyPair(id)
	if err != nil {
		return nil, err
	}
	if !model.KeyPairCanTransition(k.Status, newStatus) {
		return nil, model.NewValidationError("status", "状态流转不合法")
	}
	k.Status = newStatus
	if err := s.store.UpdateKeyPair(k); err != nil {
		return nil, err
	}
	return k, nil
}

func (s *Service) DeleteKeyPair(id string) error {
	return s.store.DeleteKeyPair(id)
}
