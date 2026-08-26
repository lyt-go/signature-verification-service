package service

import (
	"signatureservice/internal/store"
	"testing"
)

func TestRetrySuccessRejectsLateCallbackAndKeepsOneSideEffect(t *testing.T) {
	st := store.NewVerificationJobStore()
	var late func()
	effects := 0
	RunVerificationRetry(st, "job-7", func(callback func()) { late = callback }, func() { effects++ })
	late()
	job := st.Get("job-7")
	if effects != 1 {
		t.Fatalf("a retried verification should execute its external side effect once, got %d", effects)
	}
	if job.Version != 2 || job.Status != "succeeded" {
		t.Fatalf("the late first attempt must not replace the retry result, got version=%d status=%s", job.Version, job.Status)
	}
}
