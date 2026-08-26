package service

import (
	"sync"

	"signatureservice/pkg/logger"
)

type PooledVerificationAuditor struct {
	pool sync.Pool
	sink logger.AuditSink
}

func NewPooledVerificationAuditor(sink logger.AuditSink) *PooledVerificationAuditor {
	return &PooledVerificationAuditor{
		pool: sync.Pool{
			New: func() any { return &logger.AuditContext{} },
		},
		sink: sink,
	}
}

// Record 记录一次验签审计。
//
// 每条审计持有自己独占的 AuditContext：从池中取出的对象在 Emit 完成
// 之前不会归还，因此并发或重叠的审计互不影响——即使写入较慢，后一条
// 审计也不会覆盖前一条尚未读取的签名与验证者。
func (a *PooledVerificationAuditor) Record(subject, verifier string) <-chan struct{} {
	ctx := a.pool.Get().(*logger.AuditContext)
	ctx.Subject, ctx.Verifier = subject, verifier
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer a.pool.Put(ctx)
		a.sink.Emit(ctx)
	}()
	return done
}
