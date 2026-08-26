package service

import (
	"context"
	"signatureservice/internal/store"
)

func CollectVerificationBatch(ctx context.Context, inputs []string) ([]string, error) {
	records, errs := store.StreamVerificationRecords(ctx, inputs)
	result := make([]string, 0, len(inputs))
	for record := range records {
		result = append(result, record)
	}
	if err := <-errs; err != nil {
		return nil, err
	}
	return result, nil
}
