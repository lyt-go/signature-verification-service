package service

import (
	"signatureservice/internal/store"
	"testing"
)

type countingPublisher struct{ ids []string }

func (p *countingPublisher) PublishVerification(id string) error {
	p.ids = append(p.ids, id)
	return nil
}

func TestCommitRetryPublishesOneSuccessAndStoresOneRecord(t *testing.T) {
	st := store.NewVerificationTxStore(1)
	publisher := &countingPublisher{}
	if err := PersistVerifiedWithRetry(st, publisher, "verify-42"); err != nil {
		t.Fatalf("retry should recover from one temporary commit failure: %v", err)
	}
	if len(publisher.ids) != 1 {
		t.Fatalf("one committed verification should publish one success event, got %d", len(publisher.ids))
	}
	if st.RecordCount() != 1 {
		t.Fatalf("the failed transaction attempt must not leave a durable verification record, got %d", st.RecordCount())
	}
}
