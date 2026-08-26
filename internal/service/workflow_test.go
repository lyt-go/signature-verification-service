package service

import (
	"signatureservice/internal/store"
	"testing"
)

func TestDelayedExportAndCacheKeepTheFirstBatchPayload(t *testing.T) {
	cache := store.NewVerificationBatchCache()
	decoder := NewVerificationBatchDecoder(cache)
	firstGate := make(chan struct{})
	secondGate := make(chan struct{})
	close(secondGate)
	exported := make(chan string, 2)
	firstResponse := decoder.DecodeAndExport("batch-a", []byte("alpha"), firstGate, exported)
	_ = decoder.DecodeAndExport("batch-b", []byte("bravo"), secondGate, exported)
	<-exported
	close(firstGate)
	firstExport := <-exported
	if string(firstResponse) != "alpha" {
		t.Fatalf("the first response changed after decoding the second batch: %q", firstResponse)
	}
	if firstExport != "alpha" {
		t.Fatalf("the delayed first export used another batch payload: %q", firstExport)
	}
	if cached := string(cache.Get("batch-a")); cached != "alpha" {
		t.Fatalf("the cached first batch was overwritten by a later decode: %q", cached)
	}
}
