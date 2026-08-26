package service

import (
	"sort"

	"signatureservice/internal/model"
	"signatureservice/pkg/idgen"
)

func (s *Service) CreateCertificate(input model.Certificate) (*model.Certificate, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	input.ID = idgen.Hex()
	if err := s.store.CreateCertificate(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) ListCertificates(filter model.CertificateFilter, page, size int) ([]*model.Certificate, int, error) {
	all := s.store.ListCertificates()
	matched := make([]*model.Certificate, 0, len(all))
	for _, c := range all {
		if filter.Match(c) {
			matched = append(matched, c)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].ValidFrom.After(matched[j].ValidFrom)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Certificate{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetCertificate(id string) (*model.Certificate, error) {
	return s.store.GetCertificate(id)
}

func (s *Service) UpdateCertificate(id string, input model.Certificate) (*model.Certificate, error) {
	c, err := s.store.GetCertificate(id)
	if err != nil {
		return nil, err
	}
	if input.Subject != "" {
		c.Subject = input.Subject
	}
	if input.Issuer != "" {
		c.Issuer = input.Issuer
	}
	if input.SerialNumber != "" {
		c.SerialNumber = input.SerialNumber
	}
	if err := c.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateCertificate(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) UpdateCertificateStatus(id string, newStatus string) (*model.Certificate, error) {
	c, err := s.store.GetCertificate(id)
	if err != nil {
		return nil, err
	}
	if !model.CertificateCanTransition(c.Status, newStatus) {
		return nil, model.NewValidationError("status", "状态流转不合法")
	}
	c.Status = newStatus
	if err := s.store.UpdateCertificate(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) DeleteCertificate(id string) error {
	return s.store.DeleteCertificate(id)
}
