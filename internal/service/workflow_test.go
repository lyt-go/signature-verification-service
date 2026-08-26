package service

import (
	"sync"
	"testing"
	"time"

	"signatureservice/pkg/logger"
)

type blockedAuditSink struct {
	once    sync.Once
	entered chan struct{}
	release chan struct{}
	events  chan logger.AuditEvent
}

func (s *blockedAuditSink) Emit(ctx *logger.AuditContext) {
	s.once.Do(func() { close(s.entered) })
	<-s.release
	s.events <- ctx.Event()
}

func TestPooledAuditKeepsRequestIdentityAfterReturn(t *testing.T) {
	sink := &blockedAuditSink{entered: make(chan struct{}), release: make(chan struct{}), events: make(chan logger.AuditEvent, 2)}
	auditor := NewPooledVerificationAuditor(sink)
	doneA := auditor.Record("signature-a", "alice")
	select {
	case <-sink.entered:
	case <-time.After(time.Second):
		t.Fatalf("first audit consumer did not start")
	}
	doneB := auditor.Record("signature-b", "bob")
	close(sink.release)
	<-doneA
	<-doneB
	seen := map[logger.AuditEvent]int{<-sink.events: 1, <-sink.events: 1}
	if seen[(logger.AuditEvent{Subject: "signature-a", Verifier: "alice"})] != 1 || seen[(logger.AuditEvent{Subject: "signature-b", Verifier: "bob"})] != 1 {
		t.Fatalf("audit identities crossed request boundaries: %v", seen)
	}
}
