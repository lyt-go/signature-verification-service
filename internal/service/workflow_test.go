package service

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"signatureservice/internal/config"
	"signatureservice/internal/model"
	"signatureservice/internal/store"
	"signatureservice/pkg/logger"
)

func TestKeyPairSnapshotsStayImmutableDuringConcurrentUpdates(t *testing.T) {
	st := store.NewMemoryStore()
	kp := &model.KeyPair{ID: "kp-1", Name: "original", Algorithm: "rsa", KeySize: 2048, PublicKey: "public", Status: model.KeyPairActive, CreatedAt: time.Now()}
	if err := st.CreateKeyPair(kp); err != nil {
		t.Fatalf("the key pair fixture should be created: %v", err)
	}
	svc := New(st, logger.New(), config.Load())
	snapshot, err := svc.GetKeyPair("kp-1")
	if err != nil {
		t.Fatalf("the original key pair snapshot should be readable: %v", err)
	}
	var wg sync.WaitGroup
	for worker := 0; worker < 4; worker++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for n := 0; n < 80; n++ {
				_, _ = svc.UpdateKeyPair("kp-1", model.KeyPair{Name: fmt.Sprintf("worker-%d-%d", worker, n)})
			}
		}(worker)
	}
	wg.Wait()
	if snapshot.Name != "original" {
		t.Fatalf("a previously returned key pair snapshot changed after later updates: %q", snapshot.Name)
	}
}
