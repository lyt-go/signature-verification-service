package service

import (
	"signatureservice/internal/store"
	"testing"
)

func TestRecoveredBundleBuildDoesNotCachePartialState(t *testing.T) {
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("reading after a recovered build panic caused a second panic: %v", recovered)
		}
	}()
	cache := store.NewBundleCache()
	if _, err := BuildCertificateBundle(cache, "broken", true); err == nil {
		t.Fatalf("the decoder panic should be returned as a build error")
	}
	if bundle, err := cache.Get("broken"); err == nil || bundle != nil {
		t.Fatalf("a failed build must not publish a partial bundle: bundle=%v err=%v", bundle, err)
	}
	if _, err := BuildCertificateBundle(cache, "healthy", false); err != nil {
		t.Fatalf("a later healthy bundle should still build: %v", err)
	}
	if _, err := cache.Get("healthy"); err != nil {
		t.Fatalf("the healthy bundle should be readable: %v", err)
	}
}
