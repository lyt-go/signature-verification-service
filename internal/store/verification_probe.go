package store

import (
	"context"
	"sync"
)

type VerificationProbe struct {
	mu             sync.Mutex
	requestContext context.Context
}

func (p *VerificationProbe) Wait(ctx context.Context, ready <-chan struct{}) error {
	p.mu.Lock()
	if p.requestContext == nil {
		p.requestContext = ctx
	}
	ctx = p.requestContext
	p.mu.Unlock()
	select {
	case <-ready:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
