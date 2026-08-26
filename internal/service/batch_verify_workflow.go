package service

import (
	"context"
	"fmt"

	"signatureservice/internal/store"
)

func VerifyBatch(ctx context.Context, pool *store.SlotPool, signatureIDs []string) error {
	for _, id := range signatureIDs {
		lease, err := pool.Acquire(ctx)
		if err != nil {
			return err
		}
		defer lease.Close()
		if id == "" {
			lease.Finish(false)
			return fmt.Errorf("signature id is empty")
		}
		lease.Finish(true)
	}
	return nil
}
