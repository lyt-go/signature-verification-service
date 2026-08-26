package service

import (
	"sort"

	"signatureservice/internal/model"
	"signatureservice/pkg/idgen"
)

func (s *Service) CreateAlgorithm(input model.Algorithm) (*model.Algorithm, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	input.ID = idgen.Hex()
	if err := s.store.CreateAlgorithm(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) ListAlgorithms(filter model.AlgorithmFilter, page, size int) ([]*model.Algorithm, int, error) {
	all := s.store.ListAlgorithms()
	matched := make([]*model.Algorithm, 0, len(all))
	for _, a := range all {
		if filter.Match(a) {
			matched = append(matched, a)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Algorithm{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetAlgorithm(id string) (*model.Algorithm, error) {
	return s.store.GetAlgorithm(id)
}

func (s *Service) UpdateAlgorithm(id string, input model.Algorithm) (*model.Algorithm, error) {
	a, err := s.store.GetAlgorithm(id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		a.Name = input.Name
	}
	if input.Type != "" {
		a.Type = input.Type
	}
	if input.KeySize > 0 {
		a.KeySize = input.KeySize
	}
	a.Enabled = input.Enabled
	if err := a.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateAlgorithm(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) DeleteAlgorithm(id string) error {
	return s.store.DeleteAlgorithm(id)
}
