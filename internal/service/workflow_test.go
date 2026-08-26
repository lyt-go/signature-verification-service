package service

import (
	"context"
	"signatureservice/internal/store"
	"testing"
	"time"
)

func TestCancellationStopsWorkAndDoesNotPolluteNextRequest(t *testing.T) {
	probe := &store.VerificationProbe{}
	ready := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	result := make(chan error, 1)
	go func() { result <- AwaitVerification(ctx, probe, ready) }()
	cancelIgnored := false
	select {
	case err := <-result:
		if err == nil {
			cancelIgnored = true
		}
	case <-time.After(100 * time.Millisecond):
		cancelIgnored = true
		close(ready)
		<-result
	}

	directProbe := &store.VerificationProbe{}
	expired, stop := context.WithCancel(context.Background())
	stop()
	closed := make(chan struct{})
	close(closed)
	_ = AwaitDirectVerification(expired, directProbe, make(chan struct{}))
	nextErr := AwaitDirectVerification(context.Background(), directProbe, closed)
	if cancelIgnored || nextErr != nil {
		t.Fatalf("cancellation must stop the first request and stay isolated from the next request: cancelIgnored=%v nextErr=%v", cancelIgnored, nextErr)
	}
}
