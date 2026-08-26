package service

import (
	"context"
	"testing"
	"time"
)

func TestRejectedRecordTerminatesPipeline(t *testing.T) {
	done := make(chan error, 1)
	go func() {
		_, err := CollectVerificationBatch(context.Background(), []string{"sig-a", "", "sig-c"})
		done <- err
	}()
	select {
	case err := <-done:
		if err == nil {
			t.Fatalf("the invalid middle record should return its validation error")
		}
	case <-time.After(150 * time.Millisecond):
		t.Fatalf("the verification batch stayed blocked after its producer rejected a record")
	}
}
