package service

import (
	"sort"

	"signatureservice/internal/model"
	"signatureservice/pkg/idgen"
)

func (s *Service) CreateVerifyRecord(input model.VerifyRecord) (*model.VerifyRecord, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	input.ID = idgen.Hex()
	if err := s.store.CreateVerifyRecord(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) ListVerifyRecords(filter model.VerifyRecordFilter, page, size int) ([]*model.VerifyRecord, int, error) {
	all := s.store.ListVerifyRecords()
	matched := make([]*model.VerifyRecord, 0, len(all))
	for _, v := range all {
		if filter.Match(v) {
			matched = append(matched, v)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].VerifiedAt.After(matched[j].VerifiedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.VerifyRecord{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetVerifyRecord(id string) (*model.VerifyRecord, error) {
	return s.store.GetVerifyRecord(id)
}

func (s *Service) UpdateVerifyRecord(id string, input model.VerifyRecord) (*model.VerifyRecord, error) {
	vr, err := s.store.GetVerifyRecord(id)
	if err != nil {
		return nil, err
	}
	if input.Verifier != "" {
		vr.Verifier = input.Verifier
	}
	if input.PayloadHash != "" {
		vr.PayloadHash = input.PayloadHash
	}
	if err := vr.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateVerifyRecord(vr); err != nil {
		return nil, err
	}
	return vr, nil
}

func (s *Service) DeleteVerifyRecord(id string) error {
	return s.store.DeleteVerifyRecord(id)
}

func (s *Service) BatchCreateVerifyRecords(inputs []model.VerifyRecord) error {
	records := make([]*model.VerifyRecord, 0, len(inputs))
	for _, input := range inputs {
		if err := input.Validate(); err != nil {
			return err
		}
		input.ID = idgen.Hex()
		records = append(records, &input)
	}
	return s.store.BatchCreateVerifyRecords(records)
}
