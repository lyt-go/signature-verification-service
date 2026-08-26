package service

import (
	"context"
	"slices"
	"testing"
	"time"

	"signatureservice/internal/store"
)

func TestVerificationBatchReleasesSlotsAndAuditsFailure(t *testing.T) {
	pool := store.NewSlotPool(2)
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()
	if err := VerifyBatch(ctx, pool, []string{"sig-a", "sig-b", "sig-c"}); err != nil {
		t.Fatalf("three valid signatures should finish without exhausting verification slots: %v", err)
	}
	if err := VerifyBatch(context.Background(), pool, []string{"sig-d", ""}); err == nil {
		t.Fatalf("an empty signature id should be reported as a failed verification")
	}
	audits := pool.Audits()
	if len(audits) == 0 || !slices.Contains([]string{"failed", "failure"}, audits[len(audits)-1]) {
		t.Fatalf("the rejected signature should leave a failed audit, got %v", audits)
	}
}
