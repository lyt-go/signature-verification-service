package store

import (
	"context"
	"sync"
)

type SlotPool struct {
	sem    chan struct{}
	mu     sync.Mutex
	audits []string
}

func NewSlotPool(capacity int) *SlotPool { return &SlotPool{sem: make(chan struct{}, capacity)} }

type VerifyLease struct {
	pool    *SlotPool
	success bool
}

func (p *SlotPool) Acquire(ctx context.Context) (*VerifyLease, error) {
	select {
	case p.sem <- struct{}{}:
		return &VerifyLease{pool: p}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (l *VerifyLease) Finish(success bool) { l.success = success }

func (l *VerifyLease) Close() {
	<-l.pool.sem
	l.pool.mu.Lock()
	l.pool.audits = append(l.pool.audits, "success")
	l.pool.mu.Unlock()
}

func (p *SlotPool) Audits() []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]string(nil), p.audits...)
}
