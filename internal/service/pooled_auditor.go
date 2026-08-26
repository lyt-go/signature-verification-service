package service

import (
	"signatureservice/pkg/logger"
)

type PooledVerificationAuditor struct {
	pool chan *logger.AuditContext
	sink logger.AuditSink
}

func NewPooledVerificationAuditor(sink logger.AuditSink) *PooledVerificationAuditor {
	a := &PooledVerificationAuditor{pool: make(chan *logger.AuditContext, 1), sink: sink}
	a.pool <- &logger.AuditContext{}
	return a
}

func (a *PooledVerificationAuditor) Record(subject, verifier string) <-chan struct{} {
	ctx := <-a.pool
	ctx.Subject, ctx.Verifier = subject, verifier
	done := make(chan struct{})
	go func() { defer close(done); a.sink.Emit(ctx) }()
	a.pool <- ctx
	return done
}
